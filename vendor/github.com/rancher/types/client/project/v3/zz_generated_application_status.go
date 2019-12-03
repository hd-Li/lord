package client

const (
	ApplicationStatusType                   = "applicationStatus"
	ApplicationStatusFieldComponentResource = "componentResource"
)

type ApplicationStatus struct {
	ComponentResource []ComponentResources `json:"componentResource,omitempty" yaml:"componentResource,omitempty"`
}
