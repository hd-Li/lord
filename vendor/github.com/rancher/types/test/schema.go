package main

import (
	"github.com/rancher/norman/types"
	//m "github.com/rancher/norman/types/mapper"
	"github.com/rancher/types/apis/management.cattle.io/v3"
	"github.com/rancher/types/factory"
	//"github.com/rancher/types/mapper"
	//"k8s.io/api/core/v1"
)

var (
	Version = types.APIVersion{
		Version: "v3",
		Group:   "management.cattle.io",
		Path:    "/v3",
	}

	Schemas = factory.Schemas(&Version).
		Init(oamTypes)
)

func oamTypes(schemas *types.Schemas) *types.Schemas {
	return schemas.MustImport(&Version, v3.ApplicationConfigurationTemplate{})
}