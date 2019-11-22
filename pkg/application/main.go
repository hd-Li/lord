package main

import (
	"log"
	"os"
	
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"github.com/rancher/types/config"
	"github.com/rancher/norman/store/proxy"
	"github.com/rancher/norman/store/crd"
	"github.com/rancher/norman/types"
	projectschema "github.com/rancher/types/apis/project.cattle.io/v3/schema"
	projectclient "github.com/rancher/types/client/project/v3"
)

var (
	kubeconfig string = "/root/.kube/config" 
)

func main() {
	//restConfig, err := rest.InClusterConfig()
	restConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		log.Fatalf("Get restconfig failed: %s", err.Error())
		os.Exit(1)
	}
	
	userContext, err := config.NewUserOnlyContext(*restConfig)
}

func SetupApplicationCRD(apiContext *config.UserOnlyContext, config rest.Config) error {
	schemas := types.Schemas{}
	appclicationschema := apiContext.Schemas.Schema(&projectschema.Version, projectclient.)
	clientGetter, err := proxy.NewClientGetterFromConfig(config)
	if err != nil {
		log.Fatalln("create clientGetter error: %s", err.Error())
		return err
	}
	
	factory := &crd..Factory{ClientGetter: clientGetter}
}