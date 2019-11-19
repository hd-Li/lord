package v1alpha1

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RbacConfig implements the ClusterRbacConfig Custom Resource Definition for controlling Istio RBAC behavior.
// The ClusterRbacConfig Custom Resource is a singleton where only one ClusterRbacConfig should be created
// globally in the mesh and the namespace should be the same to other Istio components, which usually is `istio-system`.
//
// Below is an example of an `ClusterRbacConfig` resource called `istio-rbac-config` which enables Istio RBAC for all
// services in the default namespace.
//
// ```yaml
// apiVersion: "rbac.istio.io/v1alpha1"
// kind: ClusterRbacConfig
// metadata:
//   name: default
//   namespace: istio-system
// spec:
//   mode: ON_WITH_INCLUSION
//   inclusion:
//     namespaces: [ "default" ]
// ```
//
// <!-- go code generation tags
// +kubetype-gen
// +kubetype-gen:groupVersion=rbac.istio.io/v1alpha1
// +kubetype-gen:kubeType=RbacConfig
// +kubetype-gen:kubeType=ClusterRbacConfig
// +kubetype-gen:ClusterRbacConfig:tag=genclient:nonNamespaced
// +genclient
// +k8s:deepcopy-gen=true
// -->
type RbacConfig struct {
	v1.TypeMeta `json:",inline"`
	// +optional
	v1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines the implementation of this definition.
	// +optional
	Spec RbacConfigSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RbacConfigList is a collection of RbacConfigs.
type RbacConfigList struct {
	v1.TypeMeta `json:",inline"`
	// +optional
	v1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items       []RbacConfig `json:"items" protobuf:"bytes,2,rep,name=items"`
}

const (
	// Disable Istio RBAC completely, Istio RBAC policies will not be enforced.
	OFF RbacConfigMode = 0
	// Enable Istio RBAC for all services and namespaces. Note Istio RBAC is deny-by-default
	// which means all requests will be denied if it's not allowed by RBAC rules.
	ON RbacConfigMode = 1
	// Enable Istio RBAC only for services and namespaces specified in the inclusion field. Any other
	// services and namespaces not in the inclusion field will not be enforced by Istio RBAC policies.
	ON_WITH_INCLUSION RbacConfigMode = 2
	// Enable Istio RBAC for all services and namespaces except those specified in the exclusion field. Any other
	// services and namespaces not in the exclusion field will be enforced by Istio RBAC policies.
	ON_WITH_EXCLUSION RbacConfigMode = 3
)

type RbacConfigMode int32

// Target defines a list of services or namespaces.
type RbacConfigTarget struct {
	// A list of services.
	Services []string `protobuf:"bytes,1,rep,name=services,proto3" json:"services,omitempty"`
	// A list of namespaces.
	Namespaces           []string `protobuf:"bytes,2,rep,name=namespaces,proto3" json:"namespaces,omitempty"`
}

type EnforcementMode int32

const (
	// Policy in ENFORCED mode has impact on user experience.
	// Policy is in ENFORCED mode by default.
	ENFORCED EnforcementMode = 0
	// Policy in PERMISSIVE mode isn't enforced and has no impact on users.
	// RBAC engine run policies in PERMISSIVE mode and logs stats.
	PERMISSIVE EnforcementMode = 1
)

// RbacConfig implements the ClusterRbacConfig Custom Resource Definition for controlling Istio RBAC behavior.
// The ClusterRbacConfig Custom Resource is a singleton where only one ClusterRbacConfig should be created
// globally in the mesh and the namespace should be the same to other Istio components, which usually is `istio-system`.
//
// Below is an example of an `ClusterRbacConfig` resource called `istio-rbac-config` which enables Istio RBAC for all
// services in the default namespace.
//
// ```yaml
// apiVersion: "rbac.istio.io/v1alpha1"
// kind: ClusterRbacConfig
// metadata:
//   name: default
//   namespace: istio-system
// spec:
//   mode: ON_WITH_INCLUSION
//   inclusion:
//     namespaces: [ "default" ]
// ```
//
// <!-- go code generation tags
// +kubetype-gen
// +kubetype-gen:groupVersion=rbac.istio.io/v1alpha1
// +kubetype-gen:kubeType=RbacConfig
// +kubetype-gen:kubeType=ClusterRbacConfig
// +kubetype-gen:ClusterRbacConfig:tag=genclient:nonNamespaced
// +genclient
// +k8s:deepcopy-gen=true
// -->
type RbacConfigSpec struct {
	// Istio RBAC mode.
	Mode RbacConfigMode `protobuf:"varint,1,opt,name=mode,proto3,enum=istio.rbac.v1alpha1.RbacConfig_Mode" json:"mode,omitempty"`
	// A list of services or namespaces that should be enforced by Istio RBAC policies. Note: This field have
	// effect only when mode is ON_WITH_INCLUSION and will be ignored for any other modes.
	Inclusion *RbacConfigTarget `protobuf:"bytes,2,opt,name=inclusion,proto3" json:"inclusion,omitempty"`
	// A list of services or namespaces that should not be enforced by Istio RBAC policies. Note: This field have
	// effect only when mode is ON_WITH_EXCLUSION and will be ignored for any other modes.
	Exclusion *RbacConfigTarget `protobuf:"bytes,3,opt,name=exclusion,proto3" json:"exclusion,omitempty"`
	// $hide_from_docs
	// Indicates enforcement mode of the RbacConfig, in ENFORCED mode by default.
	// It's used to verify new RbacConfig work as expected before rolling to production.
	// When setting as PERMISSIVE, RBAC isn't enforced and has no impact on users.
	// RBAC engine run RbacConfig in PERMISSIVE mode and logs stats.
	// Invalid to set RbacConfig in PERMISSIVE and ServiceRoleBinding in ENFORCED mode.
	EnforcementMode      EnforcementMode `protobuf:"varint,4,opt,name=enforcement_mode,json=enforcementMode,proto3,enum=istio.rbac.v1alpha1.EnforcementMode" json:"enforcementMode,omitempty"`
}

//
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RbacConfig implements the ClusterRbacConfig Custom Resource Definition for controlling Istio RBAC behavior.
// The ClusterRbacConfig Custom Resource is a singleton where only one ClusterRbacConfig should be created
// globally in the mesh and the namespace should be the same to other Istio components, which usually is `istio-system`.
//
// Below is an example of an `ClusterRbacConfig` resource called `istio-rbac-config` which enables Istio RBAC for all
// services in the default namespace.
//
// ```yaml
// apiVersion: "rbac.istio.io/v1alpha1"
// kind: ClusterRbacConfig
// metadata:
//   name: default
//   namespace: istio-system
// spec:
//   mode: ON_WITH_INCLUSION
//   inclusion:
//     namespaces: [ "default" ]
// ```
//
// <!-- go code generation tags
// +kubetype-gen
// +kubetype-gen:groupVersion=rbac.istio.io/v1alpha1
// +kubetype-gen:kubeType=RbacConfig
// +kubetype-gen:kubeType=ClusterRbacConfig
// +kubetype-gen:ClusterRbacConfig:tag=genclient:nonNamespaced
// +genclient
// +k8s:deepcopy-gen=true
// -->
type ClusterRbacConfig struct {
	v1.TypeMeta `json:",inline"`
	// +optional
	v1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines the implementation of this definition.
	// +optional
	Spec  RbacConfigSpec  `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterRbacConfigList is a collection of ClusterRbacConfigs.
type ClusterRbacConfigList struct {
	v1.TypeMeta `json:",inline"`
	// +optional
	v1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items       []ClusterRbacConfig `json:"items" protobuf:"bytes,2,rep,name=items"`
}

//
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ServiceRole specification contains a list of access rules (permissions).
//
// <!-- go code generation tags
// +kubetype-gen
// +kubetype-gen:groupVersion=rbac.istio.io/v1alpha1
// +genclient
// +k8s:deepcopy-gen=true
// -->
type ServiceRole struct {
	v1.TypeMeta `json:",inline"`
	// +optional
	v1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines the implementation of this definition.
	// +optional
	Spec rbacv1alpha1.ServiceRole `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ServiceRoleList is a collection of ServiceRoles.
type ServiceRoleList struct {
	v1.TypeMeta `json:",inline"`
	// +optional
	v1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items       []ServiceRole `json:"items" protobuf:"bytes,2,rep,name=items"`
}

//
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ServiceRoleBinding assigns a ServiceRole to a list of subjects.
//
// <!-- go code generation tags
// +kubetype-gen
// +kubetype-gen:groupVersion=rbac.istio.io/v1alpha1
// +genclient
// +k8s:deepcopy-gen=true
// -->
type ServiceRoleBinding struct {
	v1.TypeMeta `json:",inline"`
	// +optional
	v1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines the implementation of this definition.
	// +optional
	Spec rbacv1alpha1.ServiceRoleBinding `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ServiceRoleBindingList is a collection of ServiceRoleBindings.
type ServiceRoleBindingList struct {
	v1.TypeMeta `json:",inline"`
	// +optional
	v1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items       []ServiceRoleBinding `json:"items" protobuf:"bytes,2,rep,name=items"`
}
