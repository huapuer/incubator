package actor

import (
	"../common/maybe"
	"../config"
	"../message"
	"context"
	"errors"
	"fmt"
	"time"
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
	SetState(Actor, string, interface{}, time.Duration) maybe.MaybeError
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
	Topo       int32
	blackBoard map[string]interface{}
}

func (this *commonActor) SetState(runner Actor, key string, value interface{}, expire time.Duration) (err maybe.MaybeError) {
	if this.blackBoard == nil {
		this.blackBoard = make(map[string]interface{})
	}
	if _, ok := this.blackBoard[key]; ok {
		err.Error(fmt.Errorf("state key already exists: %s", key))
		return
	}
	this.blackBoard[key] = value
	go func() {
		<-time.After(expire)
		runner.Receive(message.StateExpireMessage{key})
	}()
	err.Error(nil)
	return
}

func (this *commonActor) UnsetState(key string) (err maybe.MaybeError) {
	if this.blackBoard == nil {
		err.Error(errors.New("actor blackboard not set"))
		return
	}
	delete(this.blackBoard, key)
	err.Error(nil)
	return
}

func (this commonActor) GetState(key string) (ret maybe.MaybeEface) {
	if this.blackBoard == nil {
		ret.Error(errors.New("actor blackboard not set"))
		return
	}
	v, ok := this.blackBoard[key]
	if !ok {
		ret.Error(fmt.Errorf("state key does not exists: %s", key))
		return
	}
	ret.Value(v)
	return
}
