package router

import (
	"errors"
	"fmt"
	"../actor"
	"../common/maybe"
	"../config"
	"../message"
	"runtime"
)

const (
	spikeRouterClassName = "router.defaultRouter"
)

func init() {
	RegisterRouterPrototype(spikeRouterClassName, &defaultRouter{}).Test()
}

type defaultRouter struct {
	actors    []actor.Actor
	actorsNum int
}

func (this defaultRouter) New(cfg config.Config, args ...int32) config.IOC {
	ret := MaybeRouter{}
	if router, ok := cfg.Routers[spikeRouterClassName]; ok {
		return newDefaultRouter(router.ActorClass)
	}
	ret.Error(fmt.Errorf("no router class cfg found: %s", spikeRouterClassName))
	return ret
}

func newDefaultRouter(actorClassName string) (this MaybeRouter) {
	actors, err := actor.GetActors(actorClassName)
	err.Test()
	if len(actors) < 1 {
		this.Error(errors.New("actor num less than 1"))
		return
	}
	this.Value(&defaultRouter{
		actors:    actors,
		actorsNum: len(actors),
	})
	return
}

func (this defaultRouter) Route(msg message.Message) (err maybe.MaybeError) {
	seed := msg.GetHostId().Right()
	if seed < 0 {
		err.Error(fmt.Errorf("illegal hash seed: %d", seed))
		return
	}

	this.actors[seed%int64(this.actorsNum)].Receive(msg).Test()

	runtime.Gosched()

	return
}
