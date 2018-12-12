package interfaces

import (
	"github.com/incubator/common/maybe"
)

type Host interface {
	IOC
	DenseTableElement

	GetId() int64
	SetId(int64)
	Receive(RemoteMessage) maybe.MaybeError
	IsHealth() bool
	SetIP(string)
	SetPort(int)
	Start() maybe.MaybeError
}

type MaybeHost struct {
	IOC

	maybe.MaybeError
	value Host
}

func (this MaybeHost) Value(value Host) {
	this.Error(nil)
	this.value = value
}

func (this MaybeHost) Right() Host {
	this.Test()
	return this.value
}

func (this MaybeHost) New(attr interface{}, cfg Config) IOC {
	panic("not implemented.")
}
