package main

import (
	"github.com/rancher/norman/types"
	m "github.com/rancher/norman/types/mapper"
	"github.com/rancher/types/apis/management.cattle.io/v3"
	"github.com/rancher/types/factory"
	"github.com/rancher/types/mapper"
	//"k8s.io/api/core/v1"
)

var (
	Version = types.APIVersion{
		Version: "v3",
		Group:   "management.cattle.io",
		Path:    "/v3",
	}
	//Schemas = factory.Schemas(&Version)
	Schemas = factory.Schemas(&Version).Init(clusterTypes)
)

func clusterOnlyType(schemas *types.Schemas) *types.Schemas {
	return schemas.MustImport(&Version, v3.Cluster{})
}

func clusterTypes(schemas *types.Schemas) *types.Schemas {
	return schemas.
		AddMapperForType(&Version, v3.Cluster{},
			&m.Embed{Field: "status"},
			mapper.NewDropFromSchema("genericEngineConfig"),
			mapper.NewDropFromSchema("googleKubernetesEngineConfig"),
			mapper.NewDropFromSchema("azureKubernetesServiceConfig"),
			mapper.NewDropFromSchema("amazonElasticContainerServiceConfig"),
			m.DisplayName{},
		).
		AddMapperForType(&Version, v3.ClusterStatus{},
			m.Drop{Field: "serviceAccountToken"},
		).
		AddMapperForType(&Version, v3.ClusterRegistrationToken{},
			&m.Embed{Field: "status"},
		).
		AddMapperForType(&Version, v3.RancherKubernetesEngineConfig{},
			m.Drop{Field: "systemImages"},
		).
		MustImport(&Version, v3.Cluster{}).
		MustImport(&Version, v3.ClusterRegistrationToken{}).
		MustImport(&Version, v3.GenerateKubeConfigOutput{}).
		MustImport(&Version, v3.ImportClusterYamlInput{}).
		MustImport(&Version, v3.RotateCertificateInput{}).
		MustImport(&Version, v3.RotateCertificateOutput{}).
		MustImport(&Version, v3.ImportYamlOutput{}).
		MustImport(&Version, v3.ExportOutput{}).
		MustImport(&Version, v3.MonitoringInput{}).
		MustImport(&Version, v3.MonitoringOutput{}).
		MustImport(&Version, v3.RestoreFromEtcdBackupInput{}).
		MustImportAndCustomize(&Version, v3.ETCDService{}, func(schema *types.Schema) {
			schema.MustCustomizeField("extraArgs", func(field types.Field) types.Field {
				field.Default = map[string]interface{}{
					"election-timeout":   "5000",
					"heartbeat-interval": "500"}
				return field
			})
		}).
		MustImportAndCustomize(&Version, v3.Cluster{}, func(schema *types.Schema) {
			schema.MustCustomizeField("name", func(field types.Field) types.Field {
				field.Type = "dnsLabel"
				field.Nullable = true
				field.Required = false
				return field
			})
			schema.ResourceActions["generateKubeconfig"] = types.Action{
				Output: "generateKubeConfigOutput",
			}
			schema.ResourceActions["importYaml"] = types.Action{
				Input:  "importClusterYamlInput",
				Output: "importYamlOutput",
			}
			schema.ResourceActions["exportYaml"] = types.Action{
				Output: "exportOutput",
			}
			schema.ResourceActions["enableMonitoring"] = types.Action{
				Input: "monitoringInput",
			}
			schema.ResourceActions["disableMonitoring"] = types.Action{}
			schema.ResourceActions["viewMonitoring"] = types.Action{
				Output: "monitoringOutput",
			}
			schema.ResourceActions["editMonitoring"] = types.Action{
				Input: "monitoringInput",
			}
			schema.ResourceActions["backupEtcd"] = types.Action{}
			schema.ResourceActions["restoreFromEtcdBackup"] = types.Action{
				Input: "restoreFromEtcdBackupInput",
			}
			schema.ResourceActions["rotateCertificates"] = types.Action{
				Input:  "rotateCertificateInput",
				Output: "rotateCertificateOutput",
			}
		})
}