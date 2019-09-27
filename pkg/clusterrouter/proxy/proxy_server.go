package proxy

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/rancher/norman/httperror"
	"github.com/rancher/types/apis/management.cattle.io/v3"
	"github.com/rancher/types/config/dialer"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/proxy"
	"k8s.io/client-go/rest"
)

type RemoteService struct {
	cluster   *v3.Cluster
	transport transportGetter
	url       urlGetter
	auth      authGetter
}

var (
	er = &errorResponder{}
)

type urlGetter func() (url.URL, error)

type authGetter func() (string, error)

type transportGetter func() (http.RoundTripper, error)

type errorResponder struct {
}

func (e *errorResponder) Error(w http.ResponseWriter, req *http.Request, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}

func prefix(cluster *v3.Cluster) string {
	return "/k8s/clusters/" + cluster.Name
}

func New(localConfig *rest.Config, cluster *v3.Cluster, clusterLister v3.ClusterLister, factory dialer.Factory) (*RemoteService, error) {
	if cluster.Spec.Internal {
		return NewLocal(localConfig, cluster)
	}
	return NewRemote(cluster, clusterLister, factory)
}

func NewLocal(localConfig *rest.Config, cluster *v3.Cluster) (*RemoteService, error) {
	// the gvk is ignored by us, so just pass in any gvk
	hostURL, _, err := rest.DefaultServerURL(localConfig.Host, localConfig.APIPath, schema.GroupVersion{}, true)
	if err != nil {
		return nil, err
	}

	transportGetter := func() (http.RoundTripper, error) {
		return rest.TransportFor(localConfig)
	}

	rs := &RemoteService{
		cluster: cluster,
		url: func() (url.URL, error) {
			return *hostURL, nil
		},
		transport: transportGetter,
	}
	if localConfig.BearerToken != "" {
		rs.auth = func() (string, error) { return "Bearer " + localConfig.BearerToken, nil }
	} else if localConfig.Password != "" {
		rs.auth = func() (string, error) {
			return "Basic " + base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", localConfig.Username, localConfig.Password))), nil
		}
	}

	return rs, nil
}

func NewRemote(cluster *v3.Cluster, clusterLister v3.ClusterLister, factory dialer.Factory) (*RemoteService, error) {
	if !v3.ClusterConditionProvisioned.IsTrue(cluster) {
		return nil, httperror.NewAPIError(httperror.ClusterUnavailable, "cluster not provisioned")
	}

	transportGetter := func() (http.RoundTripper, error) {
		transport := &http.Transport{}

		if factory != nil {
			d, err := factory.ClusterDialer(cluster.Name)
			if err != nil {
				return nil, err
			}
			transport.Dial = d
		}
		newCluster, err := clusterLister.Get("", cluster.Name)
		if err != nil {
			return transport, err
		}
		if newCluster.Status.CACert != "" {
			certBytes, err := base64.StdEncoding.DecodeString(newCluster.Status.CACert)
			if err != nil {
				return nil, err
			}
			certs := x509.NewCertPool()
			certs.AppendCertsFromPEM(certBytes)
			transport.TLSClientConfig = &tls.Config{
				RootCAs: certs,
			}
		}
		return transport, nil
	}

	urlGetter := func() (url.URL, error) {
		newCluster, err := clusterLister.Get("", cluster.Name)
		if err != nil {
			return url.URL{}, err
		}

		u, err := url.Parse(newCluster.Status.APIEndpoint)
		if err != nil {
			return url.URL{}, err
		}
		return *u, nil
	}

	authGetter := func() (string, error) {
		newCluster, err := clusterLister.Get("", cluster.Name)
		if err != nil {
			return "", err
		}

		return "Bearer " + newCluster.Status.ServiceAccountToken, nil
	}

	return &RemoteService{
		cluster:   cluster,
		transport: transportGetter,
		url:       urlGetter,
		auth:      authGetter,
	}, nil
}

func (r *RemoteService) Close() {
}

func (r *RemoteService) Handler() http.Handler {
	return r
}

func (r *RemoteService) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	u, err := r.url()
	if err != nil {
		er.Error(rw, req, err)
		return
	}

	u.Path = strings.TrimPrefix(req.URL.Path, prefix(r.cluster))
	u.RawQuery = req.URL.RawQuery

	proto := req.Header.Get("X-Forwarded-Proto")
	if proto != "" {
		req.URL.Scheme = proto
	} else if req.TLS == nil {
		req.URL.Scheme = "http"
	} else {
		req.URL.Scheme = "https"
	}

	req.URL.Host = req.Host
	if r.auth == nil {
		req.Header.Del("Authorization")
	} else {
		token, err := r.auth()
		if err != nil {
			er.Error(rw, req, err)
			return
		}
		req.Header.Set("Authorization", token)
	}
	transport, err := r.transport()
	if err != nil {
		er.Error(rw, req, err)
		return
	}
	httpProxy := proxy.NewUpgradeAwareHandler(&u, transport, true, false, er)
	httpProxy.ServeHTTP(rw, req)
}

func (r *RemoteService) Cluster() *v3.Cluster {
	return r.cluster
}

type SimpleProxy struct {
	url                *url.URL
	transport          http.RoundTripper
	overrideHostHeader bool
}

func NewSimpleProxy(host string, caData []byte, overrideHostHeader bool) (*SimpleProxy, error) {
	hostURL, _, err := rest.DefaultServerURL(host, "", schema.GroupVersion{}, true)
	if err != nil {
		return nil, err
	}

	ht := &http.Transport{}
	if len(caData) > 0 {
		certPool := x509.NewCertPool()
		certPool.AppendCertsFromPEM(caData)
		ht.TLSClientConfig = &tls.Config{
			RootCAs: certPool,
		}
	}

	return &SimpleProxy{
		url:                hostURL,
		transport:          ht,
		overrideHostHeader: overrideHostHeader,
	}, nil
}

func (s *SimpleProxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	u := *s.url
	u.Path = req.URL.Path
	u.RawQuery = req.URL.RawQuery
	req.URL.Scheme = "https"
	req.URL.Host = req.Host
	if s.overrideHostHeader {
		req.Host = u.Host
	}
	httpProxy := proxy.NewUpgradeAwareHandler(&u, s.transport, true, false, er)
	httpProxy.ServeHTTP(rw, req)

}
