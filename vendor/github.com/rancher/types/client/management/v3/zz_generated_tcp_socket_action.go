package client

const (
	TCPSocketActionType      = "tcpSocketAction"
	TCPSocketActionFieldPort = "port"
)

type TCPSocketAction struct {
	Port int64 `json:"port,omitempty" yaml:"port,omitempty"`
}
