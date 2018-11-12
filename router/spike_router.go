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

func (this spikeRouter) New(cfg interface{}) config.IOC {
	ret := MaybeRouter{}
	if attrs, ok := cfg.(map[string]string); ok{
		if actorClass, ok := attrs["ActorClass"]; ok {
			ret.Value(newSpikeRouter(actorClass).Right())
			return ret
		}
		ret.Error(fmt.Errorf("no actor attribute found: %s", "MailBoxSize"))
		return ret
	}
	ret.Error(fmt.Errorf("illegal cfg type when new router %s", spikeRouterClassName))
	return ret
}

func newSpikeRouter(actorClassName string) (this MaybeRouter) {
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
