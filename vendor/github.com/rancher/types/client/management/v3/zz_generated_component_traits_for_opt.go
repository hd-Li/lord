package client

const (
	ComponentTraitsForOptType               = "componentTraitsForOpt"
	ComponentTraitsForOptFieldIngress       = "ingress"
	ComponentTraitsForOptFieldManualScaler  = "manualScaler"
	ComponentTraitsForOptFieldVolumeMounter = "volumeMounter"
	ComponentTraitsForOptFieldWhiteList     = "whiteList"
)

type ComponentTraitsForOpt struct {
	Ingress       *Ingress       `json:"ingress,omitempty" yaml:"ingress,omitempty"`
	ManualScaler  *ManualScaler  `json:"manualScaler,omitempty" yaml:"manualScaler,omitempty"`
	VolumeMounter *VolumeMounter `json:"volumeMounter,omitempty" yaml:"volumeMounter,omitempty"`
	WhiteList     *WhiteList     `json:"whiteList,omitempty" yaml:"whiteList,omitempty"`
}
