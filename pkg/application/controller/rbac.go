package controller

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/api/core/v1"
	"github.com/rancher/types/apis/project.cattle.io/v3"
	istiorbacv1alpha1 "github.com/lord/types/pkg/istio/apis/rbac/v1alpha1"
)

func NewServiceRoleObject(app *v3.Application, ns *corev1.Namespace) istiorbacv1alpha1.ServiceRole {
	ownerRef := GetOwnerRefFromNamespace(ns)
	
	serviceRole := istiorbacv1alpha1.ServiceRole {
		ObjectMeta: metav1.ObjectMeta{
			OwnerReferences: []metav1.OwnerReference{ownerRef},
			Namespace:       app.Namespace,
			Name:            app.Namespace + "-" + "servicerole",
		},
		Spec: istiorbacv1alpha1.ServiceRoleSpec {
			Rules: []istiorbacv1alpha1.AccessRule{
				istiorbacv1alpha1.AccessRule{
					Services: []string{"*"},
				},
			},
		},
	}
	
	return serviceRole
}

func NewServiceRoleBinding(component *v3.Component, app *v3.Application) istiorbacv1alpha1.ServiceRoleBinding {
	ownerRef := GetOwnerRef(app)
	var users map[string]string
	for _, e := range component.OptTraits.WhiteList.Users {
		users["request.auth.claims[email]"] = e
	}
	
	serviceRoleBinding := istiorbacv1alpha1.ServiceRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			OwnerReferences: []metav1.OwnerReference{ownerRef},
			Namespace:       app.Namespace,
			Name:            app.Name + "-" + component.Name + "-" + "servicerolebinding",
		},
		Spec: istiorbacv1alpha1.ServiceRoleBindingSpec {
			Subjects: []istiorbacv1alpha1.Subject{
				istiorbacv1alpha1.Subject{
					Properties: users,
				},
			},
			RoleRef: istiorbacv1alpha1.RoleRef {
				Kind: "ServiceRole",
				Name: app.Namespace + "-" + "servicerole",
			},
		},
	}
	
	return serviceRoleBinding
}