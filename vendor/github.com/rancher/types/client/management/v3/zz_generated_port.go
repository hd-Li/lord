package client

const (
	PortType               = "port"
	PortFieldContainerPort = "containerPort"
	PortFieldName          = "name"
	PortFieldProtocol      = "protocol"
)

type Port struct {
	ContainerPort int64  `json:"containerPort,omitempty" yaml:"containerPort,omitempty"`
	Name          string `json:"name,omitempty" yaml:"name,omitempty"`
	Protocol      string `json:"protocol,omitempty" yaml:"protocol,omitempty"`
}
