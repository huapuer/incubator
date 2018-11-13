package actor

import (
	"context"
	"fmt"
	"../common/maybe"
	"../config"
	"../message"
)

var (
	actorPrototype = make(map[string]Actor)
	actors         = make(map[string][]Actor)
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

type commonActor struct{
	Topo int32
}
