package client

const (
	VolumeType               = "volume"
	VolumeFieldAccessMode    = "accessMode"
	VolumeFieldDisk          = "disk"
	VolumeFieldMountPath     = "mountPath"
	VolumeFieldName          = "name"
	VolumeFieldSharingPolicy = "sharingPolicy"
)

type Volume struct {
	AccessMode    string `json:"accessMode,omitempty" yaml:"accessMode,omitempty"`
	Disk          *Disk  `json:"disk,omitempty" yaml:"disk,omitempty"`
	MountPath     string `json:"mountPath,omitempty" yaml:"mountPath,omitempty"`
	Name          string `json:"name,omitempty" yaml:"name,omitempty"`
	SharingPolicy string `json:"sharingPolicy,omitempty" yaml:"sharingPolicy,omitempty"`
}
