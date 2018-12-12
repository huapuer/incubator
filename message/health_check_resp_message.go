package message

import (
	"errors"
	"github.com/incubator/common/maybe"
	"github.com/incubator/config"
	"github.com/incubator/interfaces"
	"github.com/incubator/serialization"
	"net"
	"unsafe"
)

const (
	HealthCheckRespMessageClassName = "message.HealthCheckRespMessage"
)

func init() {
	interfaces.RegisterMessagePrototype(HealthCheckRespMessageClassName, &HealthCheckRespMessage{
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

	version int64
	health  bool

	addr string
}

func (this *HealthCheckRespMessage) Process(runner interfaces.Actor) (err maybe.MaybeError) {
	l := interfaces.GetLayer(int32(this.GetLayer())).Right()
	if this.health {
		hm, ok := l.GetTopo().LookupHost(this.GetHostId()).
			Right().(interfaces.HealthManager)
		if !ok {
			err.Error(errors.New("host is not HealthManager"))
			return
		}
		hm.Health()
	} else {
		cfg := l.GetConfig()
		cfg.(*config.Config).Layer.StartMode = config.LAYER_START_MODE_RECOVER
		msg := &PullUpMessage{
			Version: this.version,
			Cfg:     l.GetConfig().(*config.Config),
		}
		msg.SetLayer(int8(l.GetSuperLayer()))

		sl := interfaces.GetLayer(l.GetSuperLayer()).Right()
		msg.SetType(int8(sl.GetMessageType(msg).Right()))

		conn, e := net.Dial("tcp", this.addr)
		if e != nil {
			err.Error(e)
			return
		}

		_, e = conn.Write(serialization.Marshal(msg))
		if e != nil {
			err.Error(e)
			return
		}
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
