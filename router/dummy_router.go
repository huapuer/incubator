package router

import (
	"../actor"
	"../common/maybe"
	"../config"
	"../message"
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
	attrsMap, ok := attrs.(map[string]interface{})
	if !ok {
		ret.Error(fmt.Errorf("illegal cfg type when new router %s", dummyRouterClassName))
		return ret
	}
	actorSchema, ok := attrsMap["ActorSchema"]
	if !ok {
		ret.Error(fmt.Errorf("no router attribute found: %s", "ActorClass"))
		return ret
	}
	actorSchemaInt, ok := actorSchema.(int32)
	if !ok {
		ret.Error(fmt.Errorf("actor class cfg type error(expecting int): %+v", actorSchema))
		return ret
	}

	actorCfg, ok := cfg.Actors[actorSchemaInt]
	if !ok {
		ret.Error(fmt.Errorf("no actor cfg found: %s", actorSchemaInt))
		return ret
	}
	actorAttrs := actorCfg.Attributes
	if actorAttrs == nil {
		if !ok {
			ret.Error(fmt.Errorf("no actor attribute found: %d", actorSchemaInt))
			return ret
		}
	}
	newRouter := &dummyRouter{
		actor: actor.GetActorPrototype(actorCfg.Class).Right().New(actorAttrs, cfg).(actor.MaybeActor).Right(),
	}

	ret.Value(newRouter)
	return ret
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
