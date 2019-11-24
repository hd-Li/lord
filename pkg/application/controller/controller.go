package controller

import (
	"fmt"
	
	"github.com/rancher/types/config"
	"github.com/rancher/types/apis/project.cattle.io/v3/"
	"github.com/rancher/types/apis/core/v1"
	"github.com/rancher/types/apis/apps/v1beta2"
	"k8s.io/apimachinery/pkg/runtime"
)

type controller struct {
	applicationClient     v3.ApplicationInterface
	applicationLister     v3.ApplicationLister
	namespaces            v1.NamespaceInterface
	coreV1                v1.Interface
	appsV1beta2           v1beta2.Interface
}

func Register(ctx context.Context, userContext *config.UserOnlyContext) {
	c := controller{
		applicatioClient:      userContext.Project.Applications(""),
		applicatioLister:      userContext.Project.Applications("").Controller().Lister(),
		namespaces:            userContext.Core.Namespaces(""),
		coreV1:                userContext.Core,
		appsV1beta2:           userContext.Apps,
	}
	
	c.applicationClient.AddHandler(ctx, "applictionCreateOrUpdate", c.sync)
}

func (c *controller)sync(key string, app *v3.Application) (runtime.Object, error) {
	if app == nil {
		return nil, nil
	}
	
	fmt.Println(key)
	fmt.Println(*app)
	
	return nil, nil
} 