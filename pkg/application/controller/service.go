package controller

import (
	"github.com/rancher/types/apis/project.cattle.io/v3"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	istiov1alpha3 "github.com/knative/pkg/apis/istio/v1alpha3"
)

func NewServiceObject(component *v3.Component, app *v3.Application) corev1.Service {
	ownerRef := GetOwnerRef(app)
	serverPort := component.OptTraits.Ingress.ServerPort
	
	port := corev1.ServicePort {
		Port: serverPort,
		TargetPort: intstr.FromInt(int(serverPort)),
		Protocol: corev1.ProtocolTCP,
	}
	
	service := corev1.Service {
		ObjectMeta: metav1.ObjectMeta{
			OwnerReferences: []metav1.OwnerReference{ownerRef},
			Namespace:       app.Namespace,
			Name:            app.Name + "-" + component.Name + "-" + "service",
			Annotations:     map[string]string{},
		},
		Spec: corev1.ServiceSpec {
			Selector: map[string]string {
				"app": app.Name + "-" + component.Name + "-" + "workload",
			},
			Ports: []corev1.ServicePort{port},
		},
	}
	
	return service
}

func NewVirtualServiceObject(component *v3.Component, app *v3.Application) istiov1alpha3.VirtualService {
	ownerRef := GetOwnerRef(app)
	host := component.OptTraits.Ingress.Host
	service := app.Name + "-" + component.Name + "-" + "service" + "." + app.Namespace + ".svc.cluster.local"
	port := uint32(component.OptTraits.Ingress.ServerPort)
	
	virtualService := istiov1alpha3.VirtualService {
		ObjectMeta: metav1.ObjectMeta{
			OwnerReferences: []metav1.OwnerReference{ownerRef},
			Namespace:       app.Namespace,
			Name:            app.Name + "-" + component.Name + "-" + "vs",
			Annotations:     map[string]string{},
		},
		Spec: istiov1alpha3.VirtualServiceSpec {
			Gateways: []string{(app.Namespace + "-" + "gateway")},
			Hosts: []string{host},
			HTTP:  []istiov1alpha3.HTTPRoute {
				istiov1alpha3.HTTPRoute{
					Route: []istiov1alpha3.HTTPRouteDestination {
						istiov1alpha3.HTTPRouteDestination {
							Destination: istiov1alpha3.Destination{
								Host: service,
								Port: istiov1alpha3.PortSelector{
									Number: port,
								},
							},
						},
					},
				},
			},
		},
	}
	
	return virtualService
}