package interfaces

import (
	"github.com/incubator/common/maybe"
)

type Router interface {
	IOC

	Start()
	Route(RemoteMessage) maybe.MaybeError
	SimRoute(int64, int) int64
	GetActors() []Actor
	Stop()
}

type MaybeRouter struct {
	maybe.MaybeError
	value Router
}

func (this MaybeRouter) New(attrs interface{}, cfg Config) IOC {
	panic("not implemented.")
}

func (this *MaybeRouter) Value(value Router) {
	this.Error(nil)
	this.value = value
}

func (this MaybeRouter) Right() Router {
	this.Test()
	return this.value
}
