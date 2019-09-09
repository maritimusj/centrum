package model

//资源
type ResourceClass int

const (
	DefaultResClass ResourceClass = iota
	ApiResClass
	GroupResClass
	DeviceResClass
	EquipmentResClass
	MeasureResClass
	StateResClass
)

type Resource interface {
	GetResourceID() (ResourceClass, int64)
	ResourceTitle() string
	ResourceDesc() string
}
