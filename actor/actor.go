package actor

import (
	"context"
	"fmt"
	"incubator/common/maybe"
	"incubator/config"
	"incubator/message"
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

func AddActors(cfg config.Config, name string, num int) (err maybe.MaybeError) {
	if name == "" || num <= 0 {
		err.Error(fmt.Errorf("illegal actor class name or num: %s, %d", name, num))
		return
	}
	if _, ok := actors[name]; ok {
		err.Error(fmt.Errorf("actors already exists: %s", name))
		return
	}
	if prototype, ok := actorPrototype[name]; ok {
		actors[name] = make([]Actor, 0, 0)
		for i := 0; i < num; i++ {
			newActor := prototype.New(cfg).(MaybeActor).Right()
			actors[name] = append(actors[name], newActor)
		}
		return
	}
	err.Error(fmt.Errorf("actor prototype not found: %s", name))
	return
}

func GetActors(name string) (ret []Actor, err maybe.MaybeError) {
	if array, ok := actors[name]; ok {
		ret = array
		return
	}
	err.Error(fmt.Errorf("actor not found: %s", name))
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

func (this MaybeActor) New(cfg config.Config, args ...int32) config.IOC {
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
