package router

import (
	"../actor"
	"../common/maybe"
	"../config"
	"../message"
	"errors"
	"fmt"
	"context"
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

	actorSchema:=config.GetAttrInt32(attrs, "ActorSchema", config.CheckInt32GT0).Right()

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
	newRouter := &dummyRouter{
		actor: actor.GetActorPrototype(actorCfg.Class).Right().New(actorAttrs, cfg).(actor.MaybeActor).Right(),
	}

	ret.Value(newRouter)
	return ret
}

func (this dummyRouter) Start() {
	this.actor.Start(context.Background()).Test()
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
