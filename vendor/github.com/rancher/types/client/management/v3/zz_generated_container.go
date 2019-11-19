package client

const (
	ContainerType                 = "container"
	ContainerFieldArgs            = "args"
	ContainerFieldCommand         = "command"
	ContainerFieldConfig          = "config"
	ContainerFieldEnv             = "env"
	ContainerFieldImage           = "image"
	ContainerFieldImagePullPolicy = "imagePullPolicy"
	ContainerFieldImagePullSecret = "imagePullSecret"
	ContainerFieldLivenessProbe   = "livenessProbe"
	ContainerFieldName            = "name"
	ContainerFieldPorts           = "ports"
	ContainerFieldReadinessProbe  = "readinessProbe"
	ContainerFieldResources       = "resources"
	ContainerFieldSecurityContext = "securityContext"
)

type Container struct {
	Args            []string         `json:"args,omitempty" yaml:"args,omitempty"`
	Command         []string         `json:"command,omitempty" yaml:"command,omitempty"`
	Config          []ConfigFile     `json:"config,omitempty" yaml:"config,omitempty"`
	Env             []EnvVar         `json:"env,omitempty" yaml:"env,omitempty"`
	Image           string           `json:"image,omitempty" yaml:"image,omitempty"`
	ImagePullPolicy string           `json:"imagePullPolicy,omitempty" yaml:"imagePullPolicy,omitempty"`
	ImagePullSecret string           `json:"imagePullSecret,omitempty" yaml:"imagePullSecret,omitempty"`
	LivenessProbe   *HealthProbe     `json:"livenessProbe,omitempty" yaml:"livenessProbe,omitempty"`
	Name            string           `json:"name,omitempty" yaml:"name,omitempty"`
	Ports           []Port           `json:"ports,omitempty" yaml:"ports,omitempty"`
	ReadinessProbe  *HealthProbe     `json:"readinessProbe,omitempty" yaml:"readinessProbe,omitempty"`
	Resources       *CResource       `json:"resources,omitempty" yaml:"resources,omitempty"`
	SecurityContext *SecurityContext `json:"securityContext,omitempty" yaml:"securityContext,omitempty"`
}
