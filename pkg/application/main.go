package main

import (
	"log"
	"os"
	
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"github.com/rancher/types/config"
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