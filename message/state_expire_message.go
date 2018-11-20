package message

import (
	"github.com/incubator/actor"
	"github.com/incubator/common/maybe"
)

type StateExpireMessage struct {
	Key string
}

func (this *StateExpireMessage) Process(runner actor.Actor) (err maybe.MaybeError) {
	runner.UnsetState(this.Key).Test()
	err.Error(nil)
	return
}
