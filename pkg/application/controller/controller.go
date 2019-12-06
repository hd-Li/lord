package controller

import (
	"fmt"
	"context"
	"strings"
	
	"github.com/rancher/types/config"
	"github.com/rancher/types/apis/core/v1"
	"github.com/rancher/types/apis/apps/v1beta2"
	"k8s.io/apimachinery/pkg/runtime"
	"github.com/rancher/types/apis/project.cattle.io/v3"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/api/core/v1"
	istioauthnv1alpha1 "github.com/rancher/types/apis/authentication.istio.io/v1alpha1"
	istionetworkingv1alph3 "github.com/rancher/types/apis/networking.istio.io/v1alpha3"
	istiorbacv1alpha1 "github.com/rancher/types/apis/rbac.istio.io/v1alpha1"
)

var (
	LastAppliedConfigAnnotation string = "application/last-applied-configuration"
)

type controller struct {
	applicationClient     v3.ApplicationInterface
	applicationLister     v3.ApplicationLister
	namespaces            v1.NamespaceInterface
	coreV1                v1.Interface
	appsV1beta2           v1beta2.Interface
	configMapLister:      v1.ConfigMapLister
	gatewayLister:        istionetworkingv1alph3.GatewayLister
	gatewayClient:        istionetworkingv1alph3.GatewayInterface
	policyLister:         istioauthnv1alpha1.PolicyLister
	policyClient:         istioauthnv1alpha1.PolicyInterface
	clusterconfigLister:  istiorbacv1alpha1.ClusterRbacConfigLister
	clusterconfigClient:  istiorbacv1alpha1.ClusterRbacConfigInterface
}

func Register(ctx context.Context, userContext *config.UserOnlyContext) {
	c := controller{
		applicationClient:      userContext.Project.Applications(""),
		applicationLister:      userContext.Project.Applications("").Controller().Lister(),
		namespaces:            userContext.Core.Namespaces(""),
		coreV1:                userContext.Core,
		appsV1beta2:           userContext.Apps,
		configMapList:         userContext.Core.ConfigMaps("").Controller().Lister(),
		gatewayLister:         userContext.IstioNetworking.Gateways("").Controller().Lister(),
		gatewayClient:        istionetworkingv1alph3.GatewayInterface
		policyLister:         istioauthnv1alpha1.PolicyLister
		policyClient:         istioauthnv1alpha1.PolicyInterface
		clusterconfigLister:  istiorbacv1alpha1.ClusterRbacConfigLister
		clusterconfigClient:  istiorbacv1alpha1.ClusterRbacConfigInterface
	}
	
	c.applicationClient.AddHandler(ctx, "applictionCreateOrUpdate", c.sync)
}

func (c *controller)sync(key string, app *v3.Application) (runtime.Object, error) {
	if app == nil {
		return nil, nil
	}
	
	c.syncNamespaceCommon(app)
	var trusted bool = false 
	
	components := app.Spec.Components
	if len(components) == 0 {
		return nil, nil
	}
	
	if len(components[0].Containers) == 0 {
		trusted = true
	}
	
	for i, component := range components {
		if app.Status[component.Name] == nil {
			app.Status[component.Name] = v3.ComponentResources {
				ComponentId: app.Name + ":" + component.Name,
			}
		}
		
		if trusted == false {
			c.syncConfigmaps(&component, app)
			c.syncImagePullSecrets(&component, app)
			c.syncWorkload(&component, app)
		}
				
	}
	
	return nil, nil	
}

func (c *controller)syncNamespaceCommon(app *v3.Application) error {
	ns := app.Namespace
	
	
}

func (c *controller)syncConfigmaps(component *v3.Component, app *v3.Application) error {
	/*
	for _, cc := range component.Containers {
		for _, conf := range cc.Config {
			newcfgMap := GetConfigMap(conf, &cc, component, app)
			_, err := c.coreV1.ConfigMaps(configMap.Namespace).Get(configMap.Name)
			
		}
	}*/
}

func (c *controller)syncImagePullSecrets(component *v3.Component, app *v3.Application) error {
	
}

func (c *controller)syncWorkload(component *v3.Component, app *v3.Application) error {
	var f func(*v3.Component, *v3.Application) error
	
	resourceWorkloadType := "deployment"
	if resourceWorkloadType == "deployment" {
		f = NewDeployObject
	}
	
	object, err := f(component, app)
	return nil
}

func (c *controller)syncStatus (app  *v3.Application) {
}

func tmp(){
	NamespaceCommonCheck(key)
	splitted := strings.Split(key, "/")
	namespace := splitted[0]
	name := splitted[1]
	image := app.Spec.Components[0].Containers[0].Image
	ownerRef := metav1.OwnerReference{
		Name:       app.Name,
		APIVersion: app.APIVersion,
		UID:        app.UID,
		Kind:       app.Kind,
	}
	
	deploy := &appsv1beta2.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			OwnerReferences: []metav1.OwnerReference{ownerRef},
			Namespace:       namespace,
			Name:            name,
			Labels: map[string]string{
				"application": "application-test",
			},
		},
		Spec: appsv1beta2.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"application": "application-test",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"application": "application-test",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						corev1.Container{
							Name: name,
							Image: image,
						},
					},
				},
			},
		},
	}
	
	_, err := c.appsV1beta2.Deployments("").Create(deploy)
	if err != nil {
		fmt.Printf("create deploy error: %s", err.Error())
	}
	return nil, nil
}