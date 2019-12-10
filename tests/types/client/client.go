package main

import (
	"fmt"
	
	istiorabcv1alpha1 "github.com/rancher/types/apis/rbac.istio.io/v1alpha1"
	istioauthnalpha1 "github.com/rancher/types/apis/authentication.istio.io/v1alpha1"
	istionetworkingv1alph3 "github.com/rancher/types/apis/networking.istio.io/v1alpha3"
	"k8s.io/client-go/tools/clientcmd"
	"github.com/lord/types/pkg/istio/apis/rbac/v1alpha1"
	authnv1alpha1 "github.com/lord/types/pkg/istio/apis/authentication/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

var (
	kubeConfig string = "/root/.kube/config" 
)
func main() {
	restConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	
	client, err := istionetworkingv1alph3.NewForConfig(*restConfig)
	gw := client.Gateways("").Controller().Lister
	_, err = gw.Get("hd-only", "hd-only-gateway")
	fmt.Println(gw)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func rbacConfigTest(config rest.Config) error {
	istiorbac, err := istiorabcv1alpha1.NewForConfig(config)
	
	rbacConfig := &v1alpha1.ClusterRbacConfig{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:       "istio-system",
			Name:            "default",
			Labels: map[string]string{
				"application": "application-test",
			},
		},
		Spec: v1alpha1.RbacConfigSpec {
			Mode: v1alpha1.ON_WITH_INCLUSION,
			Inclusion: &v1alpha1.RbacConfigTarget {
				Namespaces: []string{"istio-test"},
			},
		},
	}
	
	_ , err = istiorbac.ClusterRbacConfigs("").Create(rbacConfig)
	
	return err
}

func policyTest(config rest.Config) error {
	authn, err := istioauthnalpha1.NewForConfig(config)
	
	originAuthenticationMethod := authnv1alpha1.OriginAuthenticationMethod {
		Jwt: &authnv1alpha1.Jwt {
			Issuer: "http://10.10.111.45:31393/auth/realms/lbj",
			JwksUri: "http://10.10.111.45:31393/auth/realms/lbj/protocol/openid-connect/certs",
		},
	}
	
	policy := &authnv1alpha1.Policy{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:       "istio-test",
			Name:            "default",
			Labels: map[string]string{
				"application": "application-test",
			},
		},
		Spec: authnv1alpha1.PolicySpec {
			Origins: []authnv1alpha1.OriginAuthenticationMethod{originAuthenticationMethod},
			PrincipalBinding: authnv1alpha1.USE_ORIGIN,
		},
	}
	
	_, err = authn.Policies("").Create(policy)
	
	return err
}