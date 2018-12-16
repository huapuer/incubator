package interfaces

import (
	"context"
	"github.com/incubator/common/maybe"
	"time"
)

type Actor interface {
	IOC

	Start(ctx context.Context) maybe.MaybeError
	Receive(Message) maybe.MaybeError
	GetState(string) maybe.MaybeEface
	UnsetState(string) maybe.MaybeError
	UnsetStateWithToken(string, int64) maybe.MaybeError
	SetState(Actor, string, interface{}, time.Duration, func(Actor)) maybe.MaybeError
	GetRouter() MaybeRouter
	SetRouter(router Router) maybe.MaybeError
	SetCancelFunc(context.CancelFunc)
	Stop()
}

type MaybeActor struct {
	maybe.MaybeError
	value Actor
}

func (this MaybeActor) New(attrs interface{}, cfg Config) IOC {
	panic("not implemented.")
}

func (this *MaybeActor) Value(value Actor) {
	this.Error(nil)
	this.value = value
}

func (this MaybeActor) Right() Actor {
	this.Test()
	return this.value
}
