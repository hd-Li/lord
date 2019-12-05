package controller

import (
	"github.com/rancher/types/apis/project.cattle.io/v3"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetDeployObject(component *v3.Component, app *v3.Application) (*appsv1beta2.Deployment, error) {
	var deploy = appsv1beta2.Deployment{}
	var objMeta = metav1.ObjectMeta{}
	var spec = appsv1beta2.DeploymentSpec{}
	component := app.
	setupDeployObjectMeta(&objMeta, app)
	deploy.ObjectMeta = objMeta
	
	selector := &metav1.LabelSelector {
		MatchLabels: map[string]string {
			"app": app.Name + "-" + 
		},
	}
	return deploy, nil
}
//application + comnentname + "kind"
func setupDeployObjectMeta(objMeta *metav1.ObjectMeta, app *v3.Application) {
	objMeta.Name = app.Name + "-" + "deployment"
	objMeta.Namespace = app.Namespace
	
	SetupOwnerForResource(objMeta, app)
	
	if app.Annotations != nil {
		objMeta.Annotations = app.Annotations
	}
	
	if app.Labels != nil {
		objMeta.Labels = app.Labels
	}
}