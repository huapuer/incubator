package router

import (
	"context"
	"errors"
	"fmt"
	"github.com/incubator/actor"
	"github.com/incubator/common/maybe"
	"github.com/incubator/config"
	"github.com/incubator/interfaces"
)

const (
	dummyRouterClassName = "router.dummyRouter"
)

func init() {
	RegisterRouterPrototype(dummyRouterClassName, &dummyRouter{}).Test()
}

type dummyRouter struct {
	actor interfaces.Actor
}

func (this dummyRouter) New(attrs interface{}, cfg interfaces.Config) interfaces.IOC {
	ret := interfaces.MaybeRouter{}

	actorSchema := config.GetAttrInt32(attrs, "ActorSchema", config.CheckInt32GT0).Right()

	actorCfg, ok := cfg.(*config.Config).ActorMap[actorSchema]
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
	newActor := actor.GetActorPrototype(actorCfg.Class).Right().New(actorAttrs, cfg).(interfaces.MaybeActor).Right()
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

////go:noescape
func (this dummyRouter) Route(msg interfaces.RemoteMessage) (err maybe.MaybeError) {
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

func (this dummyRouter) GetActors() []interfaces.Actor {
	return []interfaces.Actor{this.actor}
}

func (this dummyRouter) Stop() {
	this.actor.Stop()
}
