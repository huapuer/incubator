package message

import (
	"../actor"
	"../common/maybe"
	"../config"
	"../layer"
	"encoding/json"
	"errors"
	"fmt"
	"unsafe"
)

const (
	PullUpMessageClassName = "message.PullUpMessage"
)

func init() {
	RegisterMessagePrototype(PullUpMessageClassName, &PullUpMessage{
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

func (this *PullUpMessage) Process(runner actor.Actor) (err maybe.MaybeError) {
	if this.Cfg == nil {
		err.Error(errors.New("cfg is nil"))
		return
	}

	l := layer.GetLayer(int32(this.GetLayer())).Right()

	if this.Cfg.Layer.Recover == true {
		if l.GetVersion() != this.Version {
			err.Error(
				fmt.Errorf("layer version unmatch: origin=%d, expect=%d", l.GetVersion(), this.Version))
			return
		}
	}

	rMsg := &NodeResultMessage{}
	rMsg.SetLayer(this.GetLayer())

	rMsg.SetType(int8(l.GetMessageType(rMsg).Right()))

	rMsg.SetHostId(this.GetHostId())

	maybe.TryCatch(
		func() {
			this.Cfg.Process().Test()
			layer.GetLayer(this.Cfg.Layer.Id).Right().Start()

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
	return
}

func (this *PullUpMessage) GetSize() int32 {
	return int32(unsafe.Sizeof(*this))
}
