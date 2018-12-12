package actor

import (
	"fmt"
	"github.com/incubator/common/maybe"
	"github.com/incubator/interfaces"
	"github.com/incubator/message"
	"time"
)

type healthManager interface {
	Start(interfaces.Actor) maybe.MaybeError
}

type defaultHealthManager struct {
	heartbeatIntvl time.Duration
}

func (this defaultHealthManager) Start(runner interfaces.Actor) (err maybe.MaybeError) {
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
