package interfaces

import (
	"github.com/incubator/common/maybe"
)

type Protocal interface {
	IOC

	Pack(RemoteMessage) []byte
	Parse([]byte) (int, int)
	Decode([]byte) []byte
}

type MaybeProtocal struct {
	IOC

	maybe.MaybeError
	value Protocal
}

func (this MaybeProtocal) Value(value Protocal) {
	this.Error(nil)
	this.value = value
}

func (this MaybeProtocal) Right() Protocal {
	this.Test()
	return this.value
}

func (this MaybeProtocal) New(attrs interface{}, cfg Config) IOC {
	panic("not implemented.")
}
