package controller

import (
	"github.com/rancher/types/apis/project.cattle.io/v3"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
)

func getDeployObject(app *v3.Application) (*appsv1beta2.Deployment, error){
	var deploy appsv1beta2.Deployment{}
	
	if 
}