package main

import (
	"fmt"
	
	istiorabcv1alpha1 "github.com/rancher/types/apis/rbac.istio.io/v1alpha1"
	"k8s.io/client-go/tools/clientcmd"
	"github.com/lord/types/pkg/istio/apis/rbac/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	kubeConfig string = "/root/.kube/config" 
)
func main() {
	restConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	
	istiorbac, err := istiorabcv1alpha1.NewForConfig(*restConfig)
	
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
	if err != nil {
		fmt.Println(err.Error())
	}
}