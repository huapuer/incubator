package actor

import (
	"context"
	"errors"
	"fmt"
	"github.com/incubator/common/maybe"
	"github.com/incubator/interfaces"
)

var (
	actorPrototype = make(map[string]interfaces.Actor)
)

func RegisterActorPrototype(name string, val interfaces.Actor) (err maybe.MaybeError) {
	if _, ok := actorPrototype[name]; ok {
		err.Error(fmt.Errorf("actor prototype redefined: %s", name))
		return
	}
	actorPrototype[name] = val

	err.Error(nil)
	return
}

func GetActorPrototype(name string) (ret interfaces.MaybeActor) {
	if prototype, ok := actorPrototype[name]; ok {
		ret.Value(prototype)
		return
	}
	ret.Error(fmt.Errorf("actor prototype for class not found: %s", name))
	return
}

type commonActor struct {
	blackBoard

	Topo   int32
	r      interfaces.Router
	cancel context.CancelFunc
}

func (this commonActor) GetRouter() (ret interfaces.MaybeRouter) {
	if this.r == nil {
		ret.Error(errors.New("no router set"))
	}
	ret.Value(this.r)
	return
}

func (this *commonActor) SetRouter(r interfaces.Router) (err maybe.MaybeError) {
	if r == nil {
		err.Error(errors.New("router is nil"))
	}
	this.r = r
	err.Error(nil)
	return
}

func (this *commonActor) SetCancelFunc(cancel context.CancelFunc) {
	this.cancel = cancel
}

func (this commonActor) Stop() {
	this.cancel()
}
