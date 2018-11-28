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

	info struct {
		addr string
		cfg  *config.Config
	}
}

func (this *PullUpMessage) SetAddr(addr string) {
	this.info.addr = addr
}

func (this *PullUpMessage) SetCfg(cfg *config.Config) {
	this.info.cfg = cfg
}

func (this *PullUpMessage) Process(runner actor.Actor) (err maybe.MaybeError) {
	if this.info.cfg == nil {
		err.Error(errors.New("cfg is nil"))
		return
	}

	rMsg := layer.GetLayer(int32(this.GetLayer())).
		Right().GetMessageFromClass(NodeResultMessageClassName).
		Right().Replicate().
		Right().(*NodeResultMessage)

	rMsg.SetHostId(this.GetHostId())

	maybe.TryCatch(
		func() {
			this.info.cfg.Process().Test()
			rMsg.info.msg = "ok"
		},
		func(err error) {
			rMsg.info.msg = fmt.Sprintf("%s", err)
		})

	SendToHost(rMsg, int32(this.GetLayer()), this.GetHostId()).Test()

	return
}

func (this *PullUpMessage) GetJsonBytes() (ret maybe.MaybeBytes) {
	bytes, err := json.Marshal(this.info)
	if err != nil {
		ret.Error(err)
	} else {
		ret.Value(bytes)
	}
	return
}

func (this *PullUpMessage) SetJsonField(data []byte) (err maybe.MaybeError) {
	e := json.Unmarshal(data, this.info)
	if e != nil {
		err.Error(e)
	}
	return
}

func (this *PullUpMessage) GetSize() int32 {
	return int32(unsafe.Sizeof(*this))
}

func (this *PullUpMessage) Replicate() (ret MaybeRemoteMessage) {
	new := *this
	ret.Value(&new)
	return
}
