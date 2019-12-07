package controller

import (
	//"fmt"
	"context"
	//"strings"
	
	"github.com/rancher/types/config"
	"github.com/rancher/types/apis/core/v1"
	"github.com/rancher/types/apis/apps/v1beta2"
	"k8s.io/apimachinery/pkg/runtime"
	"github.com/rancher/types/apis/project.cattle.io/v3"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	"k8s.io/apimachinery/pkg/api/errors"
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//corev1 "k8s.io/api/core/v1"
	istioauthnv1alpha1 "github.com/rancher/types/apis/authentication.istio.io/v1alpha1"
	istionetworkingv1alph3 "github.com/rancher/types/apis/networking.istio.io/v1alpha3"
	istiorbacv1alpha1 "github.com/rancher/types/apis/rbac.istio.io/v1alpha1"
	istiov1alpha3 "github.com/knative/pkg/apis/istio/v1alpha3"
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
	configMapLister      v1.ConfigMapLister
	gatewayLister        istionetworkingv1alph3.GatewayLister
	gatewayClient        istionetworkingv1alph3.GatewayInterface
	policyLister         istioauthnv1alpha1.PolicyLister
	policyClient         istioauthnv1alpha1.PolicyInterface
	clusterconfigLister  istiorbacv1alpha1.ClusterRbacConfigLister
	clusterconfigClient  istiorbacv1alpha1.ClusterRbacConfigInterface
}

func Register(ctx context.Context, userContext *config.UserOnlyContext) {
	c := controller{
		applicationClient:      userContext.Project.Applications(""),
		applicationLister:      userContext.Project.Applications("").Controller().Lister(),
		namespaces:            userContext.Core.Namespaces(""),
		coreV1:                userContext.Core,
		appsV1beta2:           userContext.Apps,
		configMapLister:         userContext.Core.ConfigMaps("").Controller().Lister(),
		gatewayLister:         userContext.IstioNetworking.Gateways("").Controller().Lister(),
		gatewayClient:         userContext.IstioNetworking.Gateways(""),
		policyLister:          userContext.IstioAuthn.Policies("").Controller().Lister(),
		policyClient:          userContext.IstioAuthn.Policies(""),
		clusterconfigLister:   userContext.IstioRbac.ClusterRbacConfigs("").Controller().Lister(),
		clusterconfigClient:   userContext.IstioRbac.ClusterRbacConfigs(""),
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
	
	for _, component := range components {
		app.Status.ComponentResource[component.Name].ComponentId = app.Name + ":" + component.Name
		
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
	gatewayName := ns + "-" + "gateway"
	_, err := c.gatewayLister.Get(ns, gatewayName)
	if errors.IsNotFound(err) {
		gateway := NewGatewayObject(app)
		_, err = c.gatewayClient.Create(&gateway)
		if err != nil {
			return err
		}
	}
	
	_, err = c.policyLister.Get(ns, "default")
	if errors.IsNotFound(err) {
		policy := NewPolicyObject(app)
	}
	return nil
}

func (c *controller)syncConfigmaps(component *v3.Component, app *v3.Application) error {
	/*
	for _, cc := range component.Containers {
		for _, conf := range cc.Config {
			newcfgMap := GetConfigMap(conf, &cc, component, app)
			_, err := c.coreV1.ConfigMaps(configMap.Namespace).Get(configMap.Name)
			
		}
	}*/
	
	return nil
}

func (c *controller)syncImagePullSecrets(component *v3.Component, app *v3.Application) error {
	return nil
}

func (c *controller)syncWorkload(component *v3.Component, app *v3.Application) error {
	var f func(*v3.Component, *v3.Application) appsv1beta2.Deployment
	
	resourceWorkloadType := "deployment"
	if resourceWorkloadType == "deployment" {
		f = NewDeployObject
	}
	
	object := f(component, app)
	return nil
}

func (c *controller)syncStatus (app  *v3.Application) {
}