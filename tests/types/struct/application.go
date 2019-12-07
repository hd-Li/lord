package main

import (
	"fmt"
	
	"github.com/rancher/types/apis/project.cattle.io/v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main(){
	app := v3.Application{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:       "test",
			Name:            "workload",
		},
	}
	
	container := v3.ComponentContainer{
		Name: "container",
	}
	fmt.Println(app)
	fmt.Println(container)
}
