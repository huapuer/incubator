package message

import (
	"../actor"
	"../common/maybe"
	"time"
	"unsafe"
)

const (
	ActorHeartbeatMessageClassName = "message.ActorHeartbeatMessage"
)

func init() {
	RegisterMessagePrototype(ActorHeartbeatMessageClassName, &ActorHeartbeatMessage{
		commonMessage: commonMessage{
			layerId: -1,
			typ:     -1,
			master:  -1,
			hostId:  -1,
		},
	}).Test()
}

type ActorHeartbeatMessage struct {
	commonMessage

	interval time.Duration
}

func (this *ActorHeartbeatMessage) Process(runner actor.Actor) (err maybe.MaybeError) {
	runner.SetState(runner, "health_intvl", this.interval, 0, nil)
	runner.SetState(runner, "health_til", time.Now().Unix(), 0, nil)

	go func() {
		<-time.After(this.interval)
		runner.Receive(this)
	}()

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
