package message

import (
	"github.com/incubator/actor"
	"github.com/incubator/common/maybe"
)

type StateExpireMessage struct {
	Key        string
	ExpireFunc func(actor.Actor)
}

func (this StateExpireMessage) Process(runner actor.Actor) (err maybe.MaybeError) {
	runner.UnsetState(this.Key).Test()
	if this.ExpireFunc != nil {
		this.ExpireFunc(runner)
	}
	err.Error(nil)
	return
}
