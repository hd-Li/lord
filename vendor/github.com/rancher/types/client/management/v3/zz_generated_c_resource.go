package client

const (
	CResourceType         = "cResource"
	CResourceFieldCpu     = "cpu"
	CResourceFieldGpu     = "gpu"
	CResourceFieldMemory  = "memory"
	CResourceFieldVolumes = "volumes"
)

type CResource struct {
	Cpu     float64  `json:"cpu,omitempty" yaml:"cpu,omitempty"`
	Gpu     float64  `json:"gpu,omitempty" yaml:"gpu,omitempty"`
	Memory  string   `json:"memory,omitempty" yaml:"memory,omitempty"`
	Volumes []Volume `json:"volumes,omitempty" yaml:"volumes,omitempty"`
}
