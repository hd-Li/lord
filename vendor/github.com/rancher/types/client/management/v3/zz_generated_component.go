package client

const (
	ComponentType                  = "component"
	ComponentFieldArch             = "arch"
	ComponentFieldContainers       = "containers"
	ComponentFieldDevTraits        = "devTraits"
	ComponentFieldOptTraits        = "optTraits"
	ComponentFieldOsType           = "osType"
	ComponentFieldParameters       = "parameters"
	ComponentFieldWorkloadSettings = "workloadSetings"
	ComponentFieldWorkloadType     = "workloadType"
)

type Component struct {
	Arch             string                 `json:"arch,omitempty" yaml:"arch,omitempty"`
	Containers       []Container            `json:"containers,omitempty" yaml:"containers,omitempty"`
	DevTraits        *ComponentTraitsForDev `json:"devTraits,omitempty" yaml:"devTraits,omitempty"`
	OptTraits        *ComponentTraitsForOpt `json:"optTraits,omitempty" yaml:"optTraits,omitempty"`
	OsType           string                 `json:"osType,omitempty" yaml:"osType,omitempty"`
	Parameters       []Parameter            `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	WorkloadSettings []WorkloadSetting      `json:"workloadSetings,omitempty" yaml:"workloadSetings,omitempty"`
	WorkloadType     string                 `json:"workloadType,omitempty" yaml:"workloadType,omitempty"`
}
