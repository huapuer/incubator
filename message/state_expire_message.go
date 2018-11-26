package message

import (
	"../actor"
	"../common/maybe"
)

type StateExpireMessage struct {
	Key        string
	Token      int64
	ExpireFunc func(actor.Actor)
}

func (this StateExpireMessage) Process(runner actor.Actor) (err maybe.MaybeError) {
	runner.UnsetStateWithToken(this.Key, this.Token).Test()
	if this.ExpireFunc != nil {
		this.ExpireFunc(runner)
	}
	err.Error(nil)
	return
}
