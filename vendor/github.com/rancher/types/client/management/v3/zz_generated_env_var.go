package client

const (
	EnvVarType           = "envVar"
	EnvVarFieldFromParam = "fromParam"
	EnvVarFieldName      = "name"
	EnvVarFieldValue     = "value"
)

type EnvVar struct {
	FromParam string `json:"fromParam,omitempty" yaml:"fromParam,omitempty"`
	Name      string `json:"name,omitempty" yaml:"name,omitempty"`
	Value     string `json:"value,omitempty" yaml:"value,omitempty"`
}
