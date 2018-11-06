package message

import (
	"context"
	"encoding/json"
	"errors"
	"incubator/common/maybe"
	"incubator/config"
	"unsafe"
)

const (
	pullUpMessageClassName = "message.pullUpMessage"
)

func init(){
	RegisterMessagePrototype(pullUpMessageClassName, &pullUpMessage{}).Test()
}

type pullUpMessage struct {
	commonMessage
	cfg *config.Config
}

func (this *pullUpMessage) Process(ctx context.Context) (err maybe.MaybeError) {

	return
}

func (this *pullUpMessage) GetJsonBytes() (ret maybe.MaybeBytes) {
	bytes, err := json.Marshal(this)
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

func (this *pullUpMessage) GetSize() int {
	return int(unsafe.Sizeof(*this))
}

func (this *pullUpMessage) Unmarshal(data []byte) (msg MaybeMessage) {
	msg.Error(errors.New("calling abstract method:commonMessage.Unmarshal()"))
	return
}
