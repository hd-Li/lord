package client

const (
	ApplicationConfigurationSpecType            = "applicationConfigurationSpec"
	ApplicationConfigurationSpecFieldAppTraits  = "appTraits"
	ApplicationConfigurationSpecFieldComponents = "components"
	ApplicationConfigurationSpecFieldParameters = "parameters"
)

type ApplicationConfigurationSpec struct {
	AppTraits  *AppTraits  `json:"appTraits,omitempty" yaml:"appTraits,omitempty"`
	Components []Component `json:"components,omitempty" yaml:"components,omitempty"`
	Parameters []Parameter `json:"parameters,omitempty" yaml:"parameters,omitempty"`
}
