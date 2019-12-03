package main

import (
	"github.com/rancher/norman/types"
	//m "github.com/rancher/norman/types/mapper"
	mgtv3 "github.com/rancher/types/apis/management.cattle.io/v3"
	mgtschema "github.com/rancher/types/apis/management.cattle.io/v3/schema"
	v3 "github.com/rancher/types/apis/project.cattle.io/v3"
	"github.com/rancher/types/apis/project.cattle.io/v3/schema"
	"github.com/rancher/types/factory"
	//"github.com/rancher/types/mapper"
	//"k8s.io/api/core/v1"
)

var (
	Schemas = factory.Schemas(&mgtschema.Version).
		Init(oamTypes)
)

func applicationTypes(schemas *types.Schemas) *types.Schemas {
	return schemas.
		MustImport(&schema.Version, v3.Port{}).
		MustImport(&schema.Version, v3.Container{}).
		MustImport(&schema.Version, v3.Component{}).
        MustImport(&schema.Version, v3.Application{}, projectOverride{})
}

type projectOverride struct {
	types.Namespaced
	ProjectID string `norman:"type=reference[/v3/schemas/project],noupdate"`
}


func oamTypes(schemas *types.Schemas) *types.Schemas {
	return schemas.
		MustImport(&mgtschema.Version, mgtv3.ApplicationConfigurationTemplate{}).
		MustImport(&mgtschema.Version, v3.Component{}).
		MustImport(&mgtschema.Version, v3.Container{}).
		MustImport(&mgtschema.Version, v3.Port{})
}