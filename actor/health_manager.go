package actor

import (
	"../common/maybe"
	"../message"
	"fmt"
	"time"
)

type healthManager interface {
	Start(Actor) maybe.MaybeError
}

type defaultHealthManager struct {
	heartbeatIntvl time.Duration
}

func (this defaultHealthManager) Start(runner Actor) (err maybe.MaybeError) {
	if this.heartbeatIntvl <= 0 {
		err.Error(fmt.Errorf("illegal heartbeat interval: %d", this.heartbeatIntvl))
		return
	}

	msg := &message.ActorHeartbeatMessage{
		Interval: this.heartbeatIntvl,
	}

	go func() {
		for {
			runner.Receive(msg)
			time.Sleep(this.heartbeatIntvl * time.Millisecond)
		}
	}()

	err.Error(nil)
	return
}
