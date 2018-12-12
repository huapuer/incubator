package message

import (
	"github.com/incubator/common/maybe"
	"github.com/incubator/interfaces"
)

type StateExpireMessage struct {
	Key        string
	Token      int64
	ExpireFunc func(interfaces.Actor)
}

func (this StateExpireMessage) Process(runner interfaces.Actor) (err maybe.MaybeError) {
	runner.UnsetStateWithToken(this.Key, this.Token).Test()
	if this.ExpireFunc != nil {
		this.ExpireFunc(runner)
	}
	err.Error(nil)
	return
}
