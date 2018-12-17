package message

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/incubator/common/maybe"
	"github.com/incubator/config"
	"github.com/incubator/interfaces"
	"unsafe"
)

const (
	PullUpMessageClassName = "message.PullUpMessage"
)

func init() {
	interfaces.RegisterMessagePrototype(PullUpMessageClassName, &PullUpMessage{
		commonMessage: commonMessage{
			layerId: -1,
			typ:     -1,
			master:  -1,
			hostId:  -1,
		},
	}).Test()
}

type PullUpMessage struct {
	commonMessage

	Version int64
	Cfg     *config.Config
}

func (this *PullUpMessage) Process(runner interfaces.Actor) (err maybe.MaybeError) {
	if this.Cfg == nil {
		err.Error(errors.New("cfg is nil"))
		return
	}

	l := interfaces.GetLayer(int32(this.GetLayer())).Right()

	switch this.Cfg.Layer.StartMode {
	case config.LAYER_START_MODE_RECOVER:
		if l.GetVersion() != this.Version {
			err.Error(
				fmt.Errorf("layer version unmatch: origin=%d, expect=%d, ignored", l.GetVersion(), this.Version))
			return
		}
	case config.LAYER_START_MODE_REBOOT:
	case config.LAYER_START_MODE_NEW:
	default:
		err.Error(fmt.Errorf("unknown layer start mode: %d", this.Cfg.Layer.StartMode))
		return
	}

	rMsg := &NodeResultMessage{}
	rMsg.SetLayer(this.GetLayer())

	rMsg.SetType(int8(l.GetMessageType(rMsg).Right()))

	rMsg.SetHostId(this.GetHostId())

	maybe.TryCatch(
		func() {
			this.Cfg.Process().Test()
			interfaces.GetLayer(this.Cfg.Layer.Id).Right().Start()

			rMsg.info.msg = "ok"
		},
		func(err error) {
			rMsg.info.msg = fmt.Sprintf("%s", err)
		})

	SendToHost(rMsg).Test()

	return
}

func (this *PullUpMessage) GetJsonBytes() (ret maybe.MaybeBytes) {
	bytes, err := json.Marshal(this.Cfg)
	if err != nil {
		ret.Error(err)
	} else {
		ret.Value(bytes)
	}
	return
}

func (this *PullUpMessage) SetJsonField(data []byte) (err maybe.MaybeError) {
	e := json.Unmarshal(data, this.Cfg)
	if e != nil {
		err.Error(e)
	}

	err.Error(nil)
	return
}

func (this PullUpMessage) GetSize() int32 {
	return int32(unsafe.Sizeof(this))
}
