package v3

import (
	"reflect"
	
	"github.com/rancher/norman/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// PullAlways means that kubelet always attempts to pull the latest image. Container will fail If the pull fails.
	PullAlways PullPolicy = "Always"
	// PullNever means that kubelet never pulls an image, but only uses a local image. Container will fail if the image isn't present
	PullNever PullPolicy = "Never"
	// PullIfNotPresent means that kubelet pulls if the image isn't present on disk. Container will fail if the image isn't present and the pull fails.
	PullIfNotPresent PullPolicy = "IfNotPresent"
)

const (
	Server          WorkloadType = "Server"
	SingletonServer WorkloadType = "SingletonServer"
	Worker          WorkloadType = "Worker"
	SingletonWorker WorkloadType = "SingletonWorker"
	Task            WorkloadType = "Task"
	SingletonTask   WorkloadType = "SingletonTaskTask"
)

type Application struct {
	types.Namespaced
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ApplicationSpec   `json:"spec,omitempty"`
	Status ApplicationStatus `json:"status,omitempty"`
}

type ApplicationSpec struct {
	Components []Component `json:"components"`
}

type WhiteList struct {
	Users []string `json:"users,omitempty"`
}

type Ingress struct {
	Host       string `json:"host"`
	Path       string `json:"path,omitempty"`
	ServerPort int32  `json:"serverPort"`
}

type VolumeMounter struct {
	VolumeName   string `json:"volumeName"`
	StorageClass string `json:"storageClass"`
}

type ManualScaler struct {
	Replicas int32 `json:"replicas"`
}

type ComponentTraitsForOpt struct {
	ManualScaler  ManualScaler  `json:"manualScaler,omitempty"`
	VolumeMounter VolumeMounter `json:"volumeMounter,omitempty"`
	Ingress       Ingress       `json:"ingress,omitempty"`
	WhiteList     WhiteList     `json:"whiteList,omitempty"`
}

//负载均衡类型 rr;leastConn;random
//consistentType sourceIP
type IngressLB struct {
	LBType         string `json:"lbType,omitempty"`
	ConsistentType string `json:"consistentType,omitempty"`
}

type ImagePullConfig struct {
	Registry string `json:"registry,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type ComponentTraitsForDev struct {
	ImagePullConfig ImagePullConfig `json:"imagePullConfig,omitempty"`
	StaticIP        bool            `json:"staticIP,omitempty"`
	IngressLB       IngressLB       `json:"ingressLB,omitempty"`
}

type Disk struct {
	Required  string `json:"required"`
	Ephemeral bool   `json:"ephemeral"`
}

type Volume struct {
	Name          string `json:"name"`
	MountPath     string `json:"mountPath"`
	AccessMode    string `json:"accessMode,omitempty"`
	SharingPolicy string `json:"sharingPolicy,omitempty"`
	Disk          Disk   `json:"disk"`
}

type CResource struct {
	Cpu     string   `json:"cpu,omitempty"`
	Memory  string   `json:"memory,omitempty"`
	Gpu     int      `json:"gpu,omitempty"`
	Volumes []Volume `json:"volumes,omitempty"`
}

type EnvVar struct {
	Name      string `json:"name"`
	Value     string `json:"value"`
	FromParam string `json:"fromParam"`
}

type Port struct {
	Name          string `json:"name"`
	ContainerPort int32  `json:"containerPort"`
	Protocol      string `json:"protocol"`
}

type ComponentContainer struct {
	Name string `json:"name"`

	Image string `json:"image,omitempty"`

	Command []string `json:"command,omitempty"`

	Args []string `json:"args,omitempty"`

	Ports []Port `json:"ports,omitempty"`

	Env []EnvVar `json:"env,omitempty"`

	Resources CResource `json:"resources,omitempty"`

	LivenessProbe HealthProbe `json:"livenessProbe,omitempty"`

	ReadinessProbe HealthProbe `json:"readinessProbe,omitempty"`

	ImagePullPolicy PullPolicy `json:"imagePullPolicy,omitempty"`

	Config          []ConfigFile     `json:"config,omitempty"`
	ImagePullSecret string           `json:"imagePullSecret,omitempty"`
	SecurityContext *SecurityContext `json:"securityContext,omitempty"`
}

type WorkloadType string

type Component struct {
	Name       string      `json:"name"`
	Version    string      `json:"version,omitempty"`
	Parameters []Parameter `json:"parameters,omitempty"`

	WorkloadType WorkloadType `json:"workloadType"`

	OsType string `json:"osType,omitempty"`

	Arch string `json:"arch,omitempty"`

	Containers []ComponentContainer `json:"containers,omitempty"`

	WorkloadSettings []WorkloadSetting `json:"workloadSetings,omitempty"`

	DevTraits ComponentTraitsForDev `json:"devTraits,omitempty"`
	OptTraits ComponentTraitsForOpt `json:"optTraits,omitempty"`
}

//int,float,string,bool,json
type Parameter struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type"`
	Required    bool   `json:"required,omitempty"`
	Default     string `json:"default,omitempty"`
}

type SecurityContext struct{}

type ConfigFile struct {
	Path      string `json:"path"`
	FileName  string `json:"fileName"`
	Value     string `json:"value"`
	FromParam string `json:"fromParam,omitempty"`
}

type WorkloadSetting struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	Value     string `json:"value"`
	FromParam string `json:"fromParam"`
}

type HealthProbe struct {
	Exec    *ExecAction    `json:"exec,omitempty" protobuf:"bytes,1,opt,name=exec"`
	HTTPGet *HTTPGetAction `json:"httpGet,omitempty" protobuf:"bytes,2,opt,name=httpGet"`
	// TCPSocket specifies an action involving a TCP port.
	// TCP hooks not yet supported
	// TODO: implement a realistic TCP lifecycle hook
	// +optional
	TCPSocket           *TCPSocketAction `json:"tcpSocket,omitempty" protobuf:"bytes,3,opt,name=tcpSocket"`
	InitialDelaySeconds int32            `json:"initialDelaySeconds,omitempty" protobuf:"varint,2,opt,name=initialDelaySeconds"`

	TimeoutSeconds int32 `json:"timeoutSeconds,omitempty" protobuf:"varint,3,opt,name=timeoutSeconds"`

	PeriodSeconds int32 `json:"periodSeconds,omitempty" protobuf:"varint,4,opt,name=periodSeconds"`

	SuccessThreshold int32 `json:"successThreshold,omitempty" protobuf:"varint,5,opt,name=successThreshold"`

	FailureThreshold int32 `json:"failureThreshold,omitempty" protobuf:"varint,6,opt,name=failureThreshold"`
}

type TCPSocketAction struct {
	// Number or name of the port to access on the container.
	// Number must be in the range 1 to 65535.
	// Name must be an IANA_SVC_NAME.
	Port int `json:"port" protobuf:"bytes,1,opt,name=port"`
}

type HTTPGetAction struct {
	// Path to access on the HTTP server.
	// +optional
	Path string `json:"path,omitempty" protobuf:"bytes,1,opt,name=path"`
	// Name or number of the port to access on the container.
	// Number must be in the range 1 to 65535.
	// Name must be an IANA_SVC_NAME.
	Port int `json:"port" protobuf:"bytes,2,opt,name=port"`
	// Host name to connect to, defaults to the pod IP. You probably want to set
	// "Host" in httpHeaders instead.
	// +optional

	HTTPHeaders []HTTPHeader `json:"httpHeaders,omitempty" protobuf:"bytes,5,rep,name=httpHeaders"`
}

type HTTPHeader struct {
	// The header field name
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`
	// The header field value
	Value string `json:"value" protobuf:"bytes,2,opt,name=value"`
}

type ExecAction struct {
	Command []string `json:"command,omitempty" protobuf:"bytes,1,rep,name=command"`
}

type PullPolicy string

type ApplicationStatus struct {
	ComponentResource map[string]ComponentResources `json:"componentResource,omitempty"`
}

type ComponentResources struct {
	ComponentId        string   `json:"componentId,omitempty"`
	Workload           string   `json:"workload,omitempty"`
	Service            string   `json:"service,omitempty"`
	ConfigMaps         []string `json:"configMaps,omitempty"`
	ImagePullSecret    string   `json:"imagePullSecret,omitempty"`
	Gateway            string   `json:"gateway,omitempty"`
	Policy             string   `json:"policy,omitempty"`
	ClusterRbacConfig  string   `json:"clusterRbacConfig,omitempty"`
	VirtualService     string   `json:"virtualService,omitempty"`
	ServiceRole        string   `json:"serviceRole,omitempty"`
	ServiceRoleBinding string   `json:"serviceRoleBinding,omitempty"`
	DestinationRule    string   `json:"DestinationRule,omitempty"`
}
