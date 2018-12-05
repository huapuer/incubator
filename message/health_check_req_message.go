package message

import (
	"../actor"
	"../common/maybe"
	"../layer"
	"time"
	"unsafe"
)

const (
	HealthCheckReqMessageClassName = "message.HealthCheckReqMessage"
)

func init() {
	RegisterMessagePrototype(HealthCheckReqMessageClassName, &HealthCheckReqMessage{
		commonEchoMessage: commonEchoMessage{
			commonMessage: commonMessage{
				layerId: -1,
				typ:     -1,
				master:  -1,
				hostId:  -1,
			},
		},
	}).Test()
}

type HealthCheckReqMessage struct {
	commonEchoMessage
}

func (this *HealthCheckReqMessage) Process(runner actor.Actor) (err maybe.MaybeError) {

	router := runner.GetRouter().Right()

	health := true
	actors := router.GetActors()
	for _, actor := range actors {
		healthTil := actor.GetState("health_til").Right().(int64)
		healthIntvl := actor.GetState("health_intvl").Right().(int64)
		health = time.Now().Unix()-healthTil <= healthIntvl
	}

	rMsg := &HealthCheckRespMessage{
		health: health,
	}

	l := layer.GetLayer(int32(this.GetLayer())).Right()
	rMsg.SetType(int8(l.GetMessageType(rMsg).Right()))
	rMsg.SetLayer(this.GetSrcLayer())
	rMsg.SetHostId(this.GetSrcHostId())

	SendToHost(rMsg)

	return
}

func (this *HealthCheckReqMessage) GetJsonBytes() (ret maybe.MaybeBytes) {
	ret.Error(nil)
	return
}

func (this *HealthCheckReqMessage) SetJsonField(data []byte) (err maybe.MaybeError) {
	err.Error(nil)
	return
}

func (this *HealthCheckReqMessage) GetSize() int32 {
	return int32(unsafe.Sizeof(*this))
}
