package controller

import (
	"log"
	"context"
	
	
	//typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"github.com/rancher/types/config"
	"github.com/rancher/types/apis/core/v1"
	"github.com/rancher/types/apis/apps/v1beta2"
	"k8s.io/apimachinery/pkg/runtime"
	//utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"github.com/rancher/types/apis/project.cattle.io/v3"
	//appsv1beta2 "k8s.io/api/apps/v1beta2"
	"k8s.io/apimachinery/pkg/api/errors"
	//"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/record"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/api/core/v1"
	istioauthnv1alpha1 "github.com/rancher/types/apis/authentication.istio.io/v1alpha1"
	istionetworkingv1alph3 "github.com/rancher/types/apis/networking.istio.io/v1alpha3"
	istiorbacv1alpha1 "github.com/rancher/types/apis/rbac.istio.io/v1alpha1"
	//istiov1alpha3 "github.com/knative/pkg/apis/istio/v1alpha3"
)

var (
	LastAppliedConfigAnnotation string = "application/last-applied-configuration"
)

type controller struct {
	applicationClient     v3.ApplicationInterface
	applicationLister     v3.ApplicationLister
	nsClient              v1.NamespaceInterface
	coreV1                v1.Interface
	appsV1beta2           v1beta2.Interface
	deploymentLister      v1beta2.DeploymentLister
	deploymentClient      v1beta2.DeploymentInterface
	serviceLister         v1.ServiceLister
	serviceClient         v1.ServiceInterface
	virtualServiceLister  istionetworkingv1alph3.VirtualServiceLister
	virtualServiceClient  istionetworkingv1alph3.VirtualServiceInterface
	destLister            istionetworkingv1alph3.DestinationRuleLister
	destClient            istionetworkingv1alph3.DestinationRuleInterface  
	configMapLister      v1.ConfigMapLister
	gatewayLister        istionetworkingv1alph3.GatewayLister
	gatewayClient        istionetworkingv1alph3.GatewayInterface
	policyLister         istioauthnv1alpha1.PolicyLister
	policyClient         istioauthnv1alpha1.PolicyInterface
	clusterconfigLister  istiorbacv1alpha1.ClusterRbacConfigLister
	clusterconfigClient  istiorbacv1alpha1.ClusterRbacConfigInterface
	serviceRoleLister    istiorbacv1alpha1.ServiceRoleLister
	serviceRoleClient    istiorbacv1alpha1.ServiceRoleInterface
	serviceRoleBindingLister istiorbacv1alpha1.ServiceRoleBindingLister
	serviceRoleBindingClient istiorbacv1alpha1.ServiceRoleBindingInterface
	
	recorder    record.EventRecorder
}

func Register(ctx context.Context, userContext *config.UserOnlyContext) {
	/*
	utilruntime.Must(v3.AddToScheme(scheme.Scheme))
	log.Println("Creating event broadcaster")
	eventBroadcaster := record.NewBroadcaster()
	//eventBroadcaster.StartLogging(fmt.Printf)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: userContext.Core.Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: "application-controllere"})
	*/
	c := controller {
		applicationClient:     userContext.Project.Applications(""),
		applicationLister:     userContext.Project.Applications("").Controller().Lister(),
		nsClient:              userContext.Core.Namespaces(""),
		coreV1:                userContext.Core,
		appsV1beta2:           userContext.Apps,
		deploymentLister:      userContext.Apps.Deployments("").Controller().Lister(),
		deploymentClient:      userContext.Apps.Deployments(""),
		serviceLister:         userContext.Core.Services("").Controller().Lister(),
		serviceClient:         userContext.Core.Services(""),
		virtualServiceLister:   userContext.IstioNetworking.VirtualServices("").Controller().Lister(),
		virtualServiceClient:   userContext.IstioNetworking.VirtualServices(""),
		destLister:             userContext.IstioNetworking.DestinationRules("").Controller().Lister(),
		destClient:             userContext.IstioNetworking.DestinationRules(""),
		configMapLister:       userContext.Core.ConfigMaps("").Controller().Lister(),
		gatewayLister:         userContext.IstioNetworking.Gateways("").Controller().Lister(),
		gatewayClient:         userContext.IstioNetworking.Gateways(""),
		policyLister:          userContext.IstioAuthn.Policies("").Controller().Lister(),
		policyClient:          userContext.IstioAuthn.Policies(""),
		clusterconfigLister:   userContext.IstioRbac.ClusterRbacConfigs("").Controller().Lister(),
		clusterconfigClient:   userContext.IstioRbac.ClusterRbacConfigs(""),
		serviceRoleLister:     userContext.IstioRbac.ServiceRoles("").Controller().Lister(),
		serviceRoleClient:     userContext.IstioRbac.ServiceRoles(""),
		serviceRoleBindingLister: userContext.IstioRbac.ServiceRoleBindings("").Controller().Lister(),
		serviceRoleBindingClient: userContext.IstioRbac.ServiceRoleBindings(""),
	}
	
	c.applicationClient.AddHandler(ctx, "applictionCreateOrUpdate", c.sync)
}

func (c *controller)sync(key string, application *v3.Application) (runtime.Object, error) {
	if application == nil {
		return nil, nil
	}
	
	app := application.DeepCopy()
	
	c.syncNamespaceCommon(app)
	
	//the deployed app is trusted or not
	var trusted bool = false 
	
	components := app.Spec.Components
	if len(components) == 0 {
		return nil, nil
	}
	
	//if containers is nil, the app is trusted, this controller does not manage its workload's lifecycle
	if len(components[0].Containers) == 0 {
		trusted = true
	}
	
	if app.Status.ComponentResource == nil {
			app.Status.ComponentResource = map[string]v3.ComponentResources{}
	}
	
	for _, component := range components {
		if _, ok := app.Status.ComponentResource[component.Name]; !ok {
			app.Status.ComponentResource[(component.Name + "-" + component.Version)] = v3.ComponentResources{
				ComponentId: app.Name + ":" + component.Name + ":" + component.Version,
			}
		}
		
		if trusted == false {
			c.syncConfigmaps(&component, app)
			c.syncImagePullSecrets(&component, app)
			c.syncWorkload(&component, app)
		}
		
		c.syncService(&component, app)
		c.syncAuthor(&component, app)
				
	}
	
	return nil, nil	
}

func (c *controller)syncNamespaceCommon(app *v3.Application) error {
	log.Printf("Sync namespaceCommon for %s\n", app.Namespace + ":" + app.Name)
	
	var ns *corev1.Namespace
	var err error
	
	for i := 0; i < 3; i++ {
		ns, err = c.nsClient.Get(app.Namespace, metav1.GetOptions{})
		if err != nil {
			log.Printf("Get namespace object error for app %s error : %s\n", (app.Namespace + ":" + app.Name), err.Error())
		}else {
			break
		}
	}
	_, err = c.gatewayLister.Get(app.Namespace, (app.Namespace + "-" + "gateway"))
	if err != nil{
		log.Printf("Get gateway error for %s error : %s\n", (app.Namespace + ":" + app.Name), err.Error())
		
		if errors.IsNotFound(err) {
			gateway := NewGatewayObject(app, ns)
			_, err = c.gatewayClient.Create(&gateway)
			if err != nil {
				log.Printf("Create gateway error for %s error : %s\n", (app.Namespace + ":" + app.Name), err.Error())
			}
		}
	}
	log.Printf("Sync gateway done for namespace %s", app.Namespace)
	
	_, err = c.policyLister.Get(app.Namespace, "default")
	if err != nil {
		log.Printf("Get policy for %s error : %s\n", (app.Namespace + ":" + app.Name), err.Error())
		if errors.IsNotFound(err) {
			policy := NewPolicyObject(app, ns)
			_, err = c.policyClient.Create(&policy)
			if err != nil {
				log.Printf("Create policy error for %s error : %s\n", (app.Namespace + ":" + app.Name), err.Error())
			}
		}
	}
	log.Printf("Sync policy done for %s", app.Namespace)
	
	cfg, err := c.clusterconfigLister.Get("", "default")
	if err != nil{
		log.Printf("Get clusterrbacconfig for %s error : %s\n", (app.Namespace + ":" + app.Name), err.Error())
		if errors.IsNotFound(err) {
			clusterConfig := NewClusterRbacConfig(app, ns)
			_, err = c.clusterconfigClient.Create(&clusterConfig)
			if err != nil {
				log.Printf("Create clusterrbacconfig error for %s error : %s\n", (app.Namespace + ":" + app.Name), err.Error())
			}
		}
	}
	
	
	if cfg != nil {
		clusterrbacconfig := cfg.DeepCopy()
		if _, ok := clusterrbacconfig.ObjectMeta.Labels[app.Namespace]; !ok {
			clusterrbacconfig.Spec.Inclusion.Namespaces = append(clusterrbacconfig.Spec.Inclusion.Namespaces, app.Namespace)
			clusterrbacconfig.ObjectMeta.Labels[app.Namespace] = "included"
			clusterrbacconfig.Namespace = "default" //avoid the client-go bug
			_, err = c.clusterconfigClient.Update(clusterrbacconfig)
			if err != nil {
				log.Printf("Update clusterrbacconfig error for %s error : %s\n", (app.Namespace + ":" + app.Name), err.Error())
			}
		}
	}
	log.Printf("Sync clusterrbacconfig done for %s", app.Namespace)
	
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
	resourceWorkloadType := "deployment"
	if resourceWorkloadType == "deployment" {
		c.syncDeployment(component, app)
	}
	
	return nil
}

func (c *controller)syncStatus (app  *v3.Application) {
}

func (c *controller)syncDeployment(component *v3.Component, app *v3.Application) error {
	log.Printf("Sync deploy for %s", app.Namespace + ":" + component.Name)
	object := NewDeployObject(component, app)
	appliedString := GetObjectApplied(object)
	object.Annotations[LastAppliedConfigAnnotation] = appliedString
	
	deploy, err := c.deploymentLister.Get(app.Namespace, app.Name + "-" + component.Name + "-" + "workload" + "-" + component.Version)
	
	if err != nil {
		log.Printf("Get deploy for %s Error : %s\n", (app.Namespace + ":" + app.Name + ":" + component.Name), err.Error())
		if errors.IsNotFound(err) {
			_, err = c.deploymentClient.Create(&object)
			if err != nil {
				log.Printf("Create deploy for %s Error : %s\n", (app.Namespace + ":" + app.Name + ":" + component.Name), err.Error())
			}
		}
	}
	
	if deploy != nil {
		if deploy.Annotations[LastAppliedConfigAnnotation] != appliedString {
			_, err = c.deploymentClient.Update(&object)
			if err != nil {
				log.Printf("Update deploy for %s Error : %s\n", (app.Namespace + ":" + app.Name + ":" + component.Name), err.Error())
			}
		}
	}
	
	log.Printf("sync deploy for %s done!", app.Namespace + ":" + app.Name + ":" + component.Name)
	
	return nil
}

func (c *controller)syncService(component *v3.Component, app *v3.Application) error {
	object := NewServiceObject(component, app)
	objectString := GetObjectApplied(object)
	object.Annotations[LastAppliedConfigAnnotation] = objectString
	
	service, err := c.serviceLister.Get(app.Namespace, app.Name + "-" + component.Name + "-" + "service")
	if err != nil {
		log.Printf("Get service for %s Error : %s\n", (app.Namespace + ":" + app.Name + ":" + component.Name), err.Error())
		if errors.IsNotFound(err) {
			_, err = c.serviceClient.Create(&object)
			if err != nil {
				log.Printf("Create service for %s Error : %s\n", (app.Namespace + ":" + app.Name), err.Error())
			}
		}
	}
	
	if service != nil {
		if service.Annotations[LastAppliedConfigAnnotation] != objectString {
			c.serviceClient.DeleteNamespaced(service.Namespace, service.Name, &metav1.DeleteOptions{})
			_, err = c.serviceClient.Create(&object)
			if err != nil {
				log.Printf("Update(Create) Service for %s Error : %s\n", (app.Namespace + ":" + app.Name), err.Error())
			}
		}
	}
	
	_, err = c.serviceRoleLister.Get(app.Namespace, app.Name + "-" + component.Name + "-" + "servicerole")
	if err != nil {
		log.Printf("Get ServiceRole for %s Error : %s\n", (app.Name + ":" + component.Name), err.Error())
		if errors.IsNotFound(err) {
			svcRoleObject := NewServiceRoleObject(component, app)
			_, err = c.serviceRoleClient.Create(&svcRoleObject)
			if err != nil {
				log.Printf("Create ServiceRole for %s Error : %s\n", (app.Name + ":" + component.Name), err.Error())
			}
		}
	}
	
	vsObject := NewVirtualServiceObject(component, app)
	vsObjectString := GetObjectApplied(vsObject)
	vsObject.Annotations[LastAppliedConfigAnnotation] = vsObjectString
	
	vs, err := c.virtualServiceLister.Get(app.Namespace, (app.Name + "-" + component.Name + "-" + "vs"))
	if err != nil {
		if errors.IsNotFound(err) {
			_, err = c.virtualServiceClient.Create(&vsObject)
			if err != nil {
				log.Printf("Create VirtualService error for %s error : %s\n", (app.Namespace + ":" + app.Name), err.Error())
			}
		}
	}
	if vs != nil {
		if vs.Annotations[LastAppliedConfigAnnotation] !=  vsObjectString{
			vsObject.ObjectMeta.ResourceVersion = vs.ObjectMeta.ResourceVersion
			_, err = c.virtualServiceClient.Update(&vsObject)
			if err != nil {
				log.Printf("Update VirtualService error for %s error : %s\n", (app.Namespace + ":" + app.Name), err.Error())
			}
		}
	}
	
	if component.DevTraits.IngressLB.ConsistentType != "" || component.DevTraits.IngressLB.LBType != "" {
		destObject := NewDestinationruleObject(component, app)
		destObjectString := GetObjectApplied(destObject)
		destObject.Annotations[LastAppliedConfigAnnotation] = destObjectString
		
		dest, err := c.destLister.Get(app.Namespace, (app.Name + "-" + component.Name + "-" + "destinationrule"))
		if err != nil {
			log.Printf("Get DestinationRule error for %s error : %s\n", (app.Namespace + ":" + app.Name + ":" + component.Name), err.Error())
			if errors.IsNotFound(err) {
				_, err = c.destClient.Create(&destObject)
				if err != nil {
					log.Printf("Create DestinationRule error for %s error : %s\n", (app.Namespace + ":" + app.Name + ":" + component.Name), err.Error())
				}
			}
		}
		
		if dest != nil {
			if dest.Annotations[LastAppliedConfigAnnotation] != destObjectString {
				destObject.ObjectMeta.ResourceVersion = dest.ObjectMeta.ResourceVersion
				_, err := c.destClient.Update(&destObject)
				if err != nil {
					log.Printf("Update DestinationRule error for %s error : %s\n", (app.Namespace + ":" + app.Name + ":" + component.Name), err.Error())
				}
			}
		}
	}
	
	return nil
}

func (c *controller)syncAuthor(component *v3.Component, app *v3.Application) error {
	object := NewServiceRoleBinding(component, app)
	objectString := GetObjectApplied(object)
	object.Annotations[LastAppliedConfigAnnotation] = objectString
	
	serviceRoleBinding, err := c.serviceRoleBindingLister.Get(app.Namespace, object.Name)
	if err != nil {
		log.Printf("Get servicerolebinding error for %s error : %s\n", (app.Namespace + ":" + app.Name + ":" + component.Name), err.Error())
		if errors.IsNotFound(err) {
			_, err = c.serviceRoleBindingClient.Create(&object)
			if err != nil {
				log.Printf("Create servicerolebinding error for %s error : %s\n", (app.Namespace + ":" + app.Name + ":" + component.Name), err.Error())
			}
		}
	}
	
	if serviceRoleBinding != nil{
		if serviceRoleBinding.Annotations[LastAppliedConfigAnnotation] != objectString {
			object.ObjectMeta.ResourceVersion = serviceRoleBinding.ObjectMeta.ResourceVersion
			_, err = c.serviceRoleBindingClient.Update(&object)
			if err != nil{
				log.Printf("Update servicerolebinding error for %s error : %s\n", (app.Namespace + ":" + app.Name + ":" + component.Name), err.Error())
			}
		}
	}
	return nil
}