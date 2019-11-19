package client

const (
	ConfigFileType           = "configFile"
	ConfigFileFieldFromParam = "fromParam"
	ConfigFileFieldPath      = "path"
	ConfigFileFieldValue     = "value"
)

type ConfigFile struct {
	FromParam string `json:"fromParam,omitempty" yaml:"fromParam,omitempty"`
	Path      string `json:"path,omitempty" yaml:"path,omitempty"`
	Value     string `json:"value,omitempty" yaml:"value,omitempty"`
}
