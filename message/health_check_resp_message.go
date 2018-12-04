package message

import (
	"incubator/actor"
	"incubator/common/maybe"
	"unsafe"
	"encoding/json"
	"incubator/layer"
	"incubator/host"
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

	totalActor int32
	actorHealth []bool
}

func (this *HealthCheckRespMessage) Process(runner actor.Actor) (err maybe.MaybeError) {
	health := true
	for _, h := range this.actorHealth {
		if !h {
			health =false
			break
		}
	}
	if health {
		layer.GetLayer(int32(this.GetLayer())).
			Right().GetTopo().LookupHost(this.GetHostId()).
			Right().(host.HealthManager).Health()
	}
	return
}

func (this *HealthCheckRespMessage) GetJsonBytes() (ret maybe.MaybeBytes) {
	bytes, err := json.Marshal(this.actorHealth)
	if err != nil {
		ret.Error(err)
	} else {
		ret.Value(bytes)
	}
	return
}

func (this *HealthCheckRespMessage) SetJsonField(data []byte) (err maybe.MaybeError) {
	e := json.Unmarshal(data, this.actorHealth)
	if e != nil {
		err.Error(e)
	}
	return
}

func (this *HealthCheckRespMessage) GetSize() int32 {
	return int32(unsafe.Sizeof(*this))
}
