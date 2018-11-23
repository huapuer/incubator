package message

import (
	"../actor"
	"../common/maybe"
	"../config"
	"encoding/json"
	"errors"
	"unsafe"
)

const (
	pullUpMessageClassName = "message.pullUpMessage"
)

func init() {
	RegisterMessagePrototype(pullUpMessageClassName, &pullUpMessage{
		commonMessage: commonMessage{
			layer:  -1,
			typ:    -1,
			master: -1,
			hostId: -1,
		},
	}).Test()
}

type pullUpMessage struct {
	commonMessage

	cfg *config.Config
}

func (this *pullUpMessage) Process(runner actor.Actor) (err maybe.MaybeError) {
	if this.cfg == nil {
		err.Error(errors.New("cfg is nil"))
		return
	}
	this.cfg.Process().Test()
	return
}

func (this *pullUpMessage) GetJsonBytes() (ret maybe.MaybeBytes) {
	bytes, err := json.Marshal(this.cfg)
	if err != nil {
		ret.Error(err)
	} else {
		ret.Value(bytes)
	}
	return
}

func (this *pullUpMessage) SetJsonField(data []byte) (err maybe.MaybeError) {
	e := json.Unmarshal(data, this.cfg)
	if e != nil {
		err.Error(e)
	}
	return
}

func (this *pullUpMessage) GetSize() int32 {
	return int32(unsafe.Sizeof(*this))
}
