package router

import (
	"../actor"
	"../common/maybe"
	"../config"
	"../message"
	"context"
	"errors"
	"fmt"
)

const (
	dummyRouterClassName = "router.dummyRouter"
)

func init() {
	RegisterRouterPrototype(dummyRouterClassName, &dummyRouter{}).Test()
}

type dummyRouter struct {
	actor actor.Actor
}

func (this dummyRouter) New(attrs interface{}, cfg config.Config) config.IOC {
	ret := MaybeRouter{}

	actorSchema := config.GetAttrInt32(attrs, "ActorSchema", config.CheckInt32GT0).Right()

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
	newRouter := &dummyRouter{}
	newActor := actor.GetActorPrototype(actorCfg.Class).Right().New(actorAttrs, cfg).(actor.MaybeActor).Right()
	newActor.SetRouter(newRouter)
	newRouter.actor = newActor

	ret.Value(newRouter)
	return ret
}

func (this dummyRouter) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	this.actor.SetCancelFunc(cancel)
	this.actor.Start(ctx).Test()
}

//go:noescape
func (this dummyRouter) Route(msg message.RemoteMessage) (err maybe.MaybeError) {
	if this.actor == nil {
		err.Error(errors.New("actor not set"))
		return
	}

	this.actor.Receive(msg).Test()

	return
}

func (this dummyRouter) SimRoute(seed int64, actorsNum int) int64 {
	return 0
}

func (this dummyRouter) GetActors() []actor.Actor {
	return []actor.Actor{this.actor}
}

func (this dummyRouter) Stop() {
	this.actor.Stop()
}
