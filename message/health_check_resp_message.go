package message

import (
	"../actor"
	"../common/maybe"
	"../host"
	"../layer"
	"unsafe"
)

const (
	HealthCheckRespMessageClassName = "message.HealthCheckRespMessage"
)

func init() {
	RegisterMessagePrototype(HealthCheckRespMessageClassName, &HealthCheckRespMessage{
		commonMessage: commonMessage{
			layerId: -1,
			typ:     -1,
			master:  -1,
			hostId:  -1,
		},
	}).Test()
}

type HealthCheckRespMessage struct {
	commonMessage

	health bool
}

func (this *HealthCheckRespMessage) Process(runner actor.Actor) (err maybe.MaybeError) {
	if this.health {
		layer.GetLayer(int32(this.GetLayer())).
			Right().GetTopo().LookupHost(this.GetHostId()).
			Right().(host.HealthManager).Health()
	} else {
		//TODO: send the recover pullup_message
		//TODO: and need the topo recovery facility to support (topo version)
	}
	err.Error(nil)
	return
}

func (this *HealthCheckRespMessage) GetJsonBytes() (ret maybe.MaybeBytes) {
	ret.Error(nil)
	return
}

func (this *HealthCheckRespMessage) SetJsonField(data []byte) (err maybe.MaybeError) {
	err.Error(nil)
	return
}

func (this *HealthCheckRespMessage) GetSize() int32 {
	return int32(unsafe.Sizeof(*this))
}
