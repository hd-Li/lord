package client

const (
	ApplicationConfigurationStatusType       = "applicationConfigurationStatus"
	ApplicationConfigurationStatusFieldMatch = "match"
)

type ApplicationConfigurationStatus struct {
	Match bool `json:"match,omitempty" yaml:"match,omitempty"`
}
