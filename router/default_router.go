package router

import (
	"../actor"
	"../common/maybe"
	"../config"
	"../message"
	"fmt"
	"context"
)

const (
	defaultRouterClassName = "router.defaultRouter"
)

func init() {
	RegisterRouterPrototype(defaultRouterClassName, &defaultRouter{}).Test()
}

type defaultRouter struct {
	actors    []actor.Actor
	actorsNum int
	shrink int
}

func (this defaultRouter) New(attrs interface{}, cfg config.Config) config.IOC {
	ret := MaybeRouter{}

	actorSchema:=config.GetAttrInt32(attrs, "ActorSchema", config.CheckInt32GT0).Right()
	actorNum:=config.GetAttrInt(attrs, "ActorNum", config.CheckIntGT0).Right()
	shrink := config.GetAttrInt(attrs, "Shrink", config.CheckIntGT0).Right()

	actorCfg, ok := cfg.Actors[actorSchema]
	if !ok {
		ret.Error(fmt.Errorf("no actor cfg found: %s", actorSchema))
		return ret
	}
	actorAttrs := actorCfg.Attributes
	if actorAttrs == nil {
		if !ok {
			ret.Error(fmt.Errorf("no actor attribute found: %d", actorSchema))
			return ret
		}
	}
	newRouter := &defaultRouter{
		actorsNum: actorNum,
		actors:    make([]actor.Actor, 0, 0),
		shrink:shrink,
	}
	for i := 0; i < actorNum; i++ {
		newActor := actor.GetActorPrototype(actorCfg.Class).Right().New(actorAttrs, cfg).(actor.MaybeActor).Right()
		newActor.SetRouter(newRouter)
		newRouter.actors = append(newRouter.actors, newActor)
	}
	ret.Value(newRouter)
	return ret
}

func (this defaultRouter) Start() {
	for _, actor := range this.actors {
		actor.Start(context.Background()).Test()
	}
}

//go:noescape
func (this defaultRouter) Route(msg message.RemoteMessage) (err maybe.MaybeError) {
	seed := msg.GetHostId()
	if seed < 0 {
		err.Error(fmt.Errorf("illegal hash seed: %d", seed))
		return
	}

	this.actors[(seed/int64(this.shrink))%int64(this.actorsNum)].Receive(msg).Test()

	return
}

func (this defaultRouter) SimRoute(seed int64, actorsNum int) int64 {
	return (seed/int64(this.shrink))%int64(actorsNum)
}

func (this defaultRouter) GetActors() []actor.Actor {
	return this.actors
}