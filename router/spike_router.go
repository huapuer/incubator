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

type spikeRouter struct {
	actors    []actor.Actor
	actorsNum int
}

func (this spikeRouter) New(attrs interface{}, cfg config.Config) config.IOC {
	ret := MaybeRouter{}
	attrsMap, ok := attrs.(map[string]interface{})
	if !ok{
		ret.Error(fmt.Errorf("illegal cfg type when new router %s", defaultRouterClassName))
		return ret
	}
	actorSchema, ok := attrsMap["ActorSchema"]
	if !ok{
		ret.Error(fmt.Errorf("no router attribute found: %s", "ActorClass"))
		return ret
	}
	actorSchemaInt, ok := actorSchema.(int32)
	if !ok {
		ret.Error(fmt.Errorf("actor class cfg type error(expecting int): %+v", actorSchema))
		return ret
	}
	actorNum, ok := attrsMap["ActorNum"]
	if !ok{
		ret.Error(fmt.Errorf("no router attribute found: %s", "ActorNum"))
		return ret
	}
	actorNumInt, ok := actorNum.(int)
	if !ok {
		ret.Error(fmt.Errorf("actor num cfg type error(expecting int): %+v", actorNumInt))
		return ret
	}

	actorCfg, ok := cfg.Actors[actorSchemaInt]
	if !ok {
		ret.Error(fmt.Errorf("no actor cfg found: %s", actorSchemaInt))
		return ret
	}
	actorAttrs := actorCfg.Attributes
	if actorAttrs == nil {
		if !ok{
			ret.Error(fmt.Errorf("no actor attribute found: %d", actorSchemaInt))
			return ret
		}
	}
	newRouter := &defaultRouter{
		actorsNum:actorNumInt,
		actors:make([]actor.Actor,0,0),
	}
	for i:=0;i<actorNumInt;i++{
		newActor := actor.GetActorPrototype(actorCfg.Class).Right().New(actorAttrs, cfg).(actor.MaybeActor).Right()
		newRouter.actors = append(newRouter.actors, newActor)
	}
	ret.Value(newRouter)
	return ret
}

func (this spikeRouter) Route(msg message.Message) (err maybe.MaybeError) {
	seed := msg.GetHostId().Right()
	if seed < 0 {
		err.Error(fmt.Errorf("illegal hash seed: %d", seed))
		return
	}

	this.actors[seed%int64(this.actorsNum)].Receive(msg).Test()

	runtime.Gosched()

	return
}
