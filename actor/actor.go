package actor

import (
	"../common/maybe"
	"../config"
	"../message"
	"context"
	"fmt"
	"time"
	"incubator/router"
	"errors"
)

var (
	actorPrototype = make(map[string]Actor)
)

func RegisterActorPrototype(name string, val Actor) (err maybe.MaybeError) {
	if _, ok := actorPrototype[name]; ok {
		err.Error(fmt.Errorf("actor prototype redefined: %s", name))
		return
	}
	actorPrototype[name] = val
	return
}

func GetActorPrototype(name string) (ret MaybeActor) {
	if prototype, ok := actorPrototype[name]; ok {
		ret.Value(prototype)
		return
	}
	ret.Error(fmt.Errorf("actor prototype for class not found: %s", name))
	return
}

type Actor interface {
	config.IOC

	Start(ctx context.Context) maybe.MaybeError
	Receive(message.Message) maybe.MaybeError
	GetState(string) maybe.MaybeEface
	UnsetState(string) maybe.MaybeError
	UnsetStateWithToken(string, int64) maybe.MaybeError
	SetState(Actor, string, interface{}, time.Duration, func(Actor)) maybe.MaybeError
	GetRouter() router.MaybeRouter
	SetRouter(router router.Router) maybe.MaybeError
}

type MaybeActor struct {
	maybe.MaybeError
	value Actor
}

func (this MaybeActor) New(attrs interface{}, cfg config.Config) config.IOC {
	panic("not implemented.")
}

func (this MaybeActor) Value(value Actor) {
	this.Error(nil)
	this.value = value
}

func (this MaybeActor) Right() Actor {
	this.Test()
	return this.value
}

type commonActor struct {
	blackBoard

	Topo int32
	r router.Router
}

func (this commonActor) GetRouter() (ret router.MaybeRouter) {
	if this.r == nil {
		ret.Error(errors.New("no router set"))
	}
	ret.Value(this.r)
	return
}

func (this *commonActor) SetRouter(r router.Router) (err maybe.MaybeError){
	if r == nil {
		err.Error(errors.New("router is nil"))
	}
	this.r = r
	err.Error(nil)
	return
}
