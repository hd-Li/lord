package client

const (
	IngressType            = "ingress"
	IngressFieldHost       = "host"
	IngressFieldPath       = "path"
	IngressFieldServerPort = "serverPort"
)

type Ingress struct {
	Host       string `json:"host,omitempty" yaml:"host,omitempty"`
	Path       string `json:"path,omitempty" yaml:"path,omitempty"`
	ServerPort int64  `json:"serverPort,omitempty" yaml:"serverPort,omitempty"`
}
