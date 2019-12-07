package controller

import (
	"github.com/rancher/types/apis/project.cattle.io/v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	istiov1alpha3 "github.com/knative/pkg/apis/istio/v1alpha3"
	istioauthnv1alphav1 "github.com/lord/types/pkg/istio/apis/authentication/v1alpha1"
)

func GetOwnerRef(app *v3.Application) metav1.OwnerReference {
	ownerRef := metav1.OwnerReference{
		Name:       app.Namespace + ":" + app.Name,
		APIVersion: app.APIVersion,
		UID:        app.UID,
		Kind:       app.Kind,
	}
	
	return ownerRef
}

func NewGatewayObject(app *v3.Application) istiov1alpha3.Gateway {
	ownerRef := GetOwnerRef(app)
	
	gateway:= istiov1alpha3.Gateway {
		ObjectMeta: metav1.ObjectMeta{
			OwnerReferences: []metav1.OwnerReference{ownerRef},
			Namespace:       app.Namespace,
			Name:            app.Namespace + "-" "gateway",
		},
		Spec: istiov1alpha3.GatewaySpec {
			Selector: map[string]string{"app": "istio-ingressgateway"},
			Servers: []istiov1alpha3.Server{
				{
					Hosts: []string{"*"},
					Port: istiov1alpha3.Port{
						Name: "http",
						Number: "80",
						Protocol: istiov1alpha3.ProtocolHTTP,
					},
				},
			}
		},
	}
	
	return gateway
}

func NewPolicyObject(app *v3.Application) istioauthnv1alphav1.Policy {
	ownerRef := GetOwnerRef(app)
	
	policy := istioauthnv1alphav1.Policy {
		ObjectMeta: metav1.ObjectMeta{
			OwnerReferences: []metav1.OwnerReference{ownerRef},
			Namespace:       app.Namespace,
			Name:            app.Namespace + "-" + "policy",
		},
		Spec: istioauthnv1alphav1.PolicySpec {
			
		},
	}
}