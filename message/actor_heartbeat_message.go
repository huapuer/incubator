package message

import (
	"github.com/incubator/common/maybe"
	"github.com/incubator/interfaces"
	"time"
	"unsafe"
)

type ActorHeartbeatMessage struct {
	Interval time.Duration
}

func (this *ActorHeartbeatMessage) Process(runner interfaces.Actor) (err maybe.MaybeError) {
	runner.SetState(runner, "health_intvl", this.Interval, 0, nil)
	runner.SetState(runner, "health_til", time.Now().Unix(), 0, nil)

	go func() {
		<-time.After(this.Interval * time.Millisecond)
		runner.Receive(this)
	}()

	err.Error(nil)
	return
}

func (this *ActorHeartbeatMessage) GetJsonBytes() (ret maybe.MaybeBytes) {
	ret.Error(nil)
	return
}

func (this *ActorHeartbeatMessage) SetJsonField(data []byte) (err maybe.MaybeError) {
	err.Error(nil)
	return
}

func (this *ActorHeartbeatMessage) GetSize() int32 {
	return int32(unsafe.Sizeof(*this))
}
