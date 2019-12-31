package controller

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/rancher/types/apis/project.cattle.io/v3"
	istiorbacv1alpha1 "github.com/lord/types/pkg/istio/apis/rbac/v1alpha1"
)

func NewServiceRoleObject(component *v3.Component, app *v3.Application) istiorbacv1alpha1.ServiceRole {
	ownerRef := GetOwnerRef(app)
	
	serviceRole := istiorbacv1alpha1.ServiceRole {
		ObjectMeta: metav1.ObjectMeta{
			OwnerReferences: []metav1.OwnerReference{ownerRef},
			Namespace:       app.Namespace,
			Name:            app.Name + "-" + component.Name + "-" + "servicerole",
			Annotations:     map[string]string{},
		},
		Spec: istiorbacv1alpha1.ServiceRoleSpec {
			Rules: []istiorbacv1alpha1.AccessRule{
				istiorbacv1alpha1.AccessRule{
					Services: []string{(app.Name+"-"+component.Name+"-"+"service"+"."+app.Namespace+".svc.cluster.local")},
				},
			},
		},
	}
	
	return serviceRole
}

func NewServiceRoleBinding(component *v3.Component, app *v3.Application) istiorbacv1alpha1.ServiceRoleBinding {
	ownerRef := GetOwnerRef(app)
	
	subjects := []istiorbacv1alpha1.Subject{}
	
	for _, e := range component.OptTraits.WhiteList.Users {
		subject := istiorbacv1alpha1.Subject{
			Properties: map[string]string{
				"request.auth.claims[email]": e,
			},
		}
		
		subjects = append(subjects, subject)
	}
	
	serviceRoleBinding := istiorbacv1alpha1.ServiceRoleBinding{
		TypeMeta: metav1.TypeMeta{
			Kind: "ServiceRoleBinding",
			APIVersion: "rbac.istio.io/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			OwnerReferences: []metav1.OwnerReference{ownerRef},
			Namespace:       app.Namespace,
			Name:            app.Name + "-" + component.Name + "-" + "servicerolebinding",
			Annotations:     map[string]string{},
		},
		Spec: istiorbacv1alpha1.ServiceRoleBindingSpec {
			Subjects: subjects,
			RoleRef: istiorbacv1alpha1.RoleRef {
				Kind: "ServiceRole",
				Name: app.Name+"-"+component.Name+"-"+"servicerole",
			},
		},
	}
	
	return serviceRoleBinding
}