package client

const (
	HTTPGetActionType             = "httpGetAction"
	HTTPGetActionFieldHTTPHeaders = "httpHeaders"
	HTTPGetActionFieldPath        = "path"
	HTTPGetActionFieldPort        = "port"
)

type HTTPGetAction struct {
	HTTPHeaders []HTTPHeader `json:"httpHeaders,omitempty" yaml:"httpHeaders,omitempty"`
	Path        string       `json:"path,omitempty" yaml:"path,omitempty"`
	Port        int64        `json:"port,omitempty" yaml:"port,omitempty"`
}
