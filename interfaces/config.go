package interfaces

import "github.com/incubator/common/maybe"

type Config interface {
	Process() maybe.MaybeError
	GetLayerId() int32
	GetLayerClass() string
	GetLayerAttr() interface{}
}
