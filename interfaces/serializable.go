package interfaces

import "github.com/incubator/common/maybe"

type Serializable interface {
	GetSize() int32
	GetJsonBytes() maybe.MaybeBytes
	SetJsonField([]byte) maybe.MaybeError
}

type MaybeSerializable struct {
	maybe.MaybeError
	value Serializable
}

func (this MaybeSerializable) Value(value Serializable) {
	this.Error(nil)
	this.value = value
}

func (this MaybeSerializable) Right() Serializable {
	this.Test()
	return this.value
}
