package controller

import (
	"github.com/rancher/types/apis/project.cattle.io/v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NameSpaceCommonCheck(ns string) {
	gateway
	policy
	clusterconfig
}

func SetupOwnerForResource(objMeta *metav1.ObjectMeta, app *v3.Application) {
	ownerRef := metav1.OwnerReference{
		Name:       app.Namespace + ":" + app.Name,
		APIVersion: app.APIVersion,
		UID:        app.UID,
		Kind:       app.Kind,
	}
	
	objMeta.OwnerReferences = append(objMeta.OwnerReferences, ownerRef)
}