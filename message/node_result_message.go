package message

import (
	"incubator/actor"
	"incubator/common/maybe"
	"encoding/json"
	"unsafe"
	"fmt"
)

const (
	NodeResultMessageClassName = "message.resultMessage"
)

func init() {
	RegisterMessagePrototype(NodeResultMessageClassName, &NodeResultMessage{
		commonSessionedMessage: commonSessionedMessage{
			commonMessage:commonMessage{
				layer:  -1,
				typ:    -1,
				master: -1,
				hostId: -1,
			},
		},
	}).Test()
}

type NodeResultMessage struct {
	commonSessionedMessage

	info struct {
		addr string
		msg string
	}
}

func (this *NodeResultMessage) Process(runner actor.Actor) (err maybe.MaybeError) {
	fmt.Printf("node: %s, msg: %s", this.info.addr, this.info.msg)
	error.Error(nil)
	return
}

func (this *NodeResultMessage) GetJsonBytes() (ret maybe.MaybeBytes) {
	bytes, err := json.Marshal(this.info)
	if err != nil {
		ret.Error(err)
	} else {
		ret.Value(bytes)
	}
	return
}

func (this *NodeResultMessage) SetJsonField(data []byte) (err maybe.MaybeError) {
	e := json.Unmarshal(data, this.info)
	if e != nil {
		err.Error(e)
	}
	return
}

func (this *NodeResultMessage) GetSize() int32 {
	return int32(unsafe.Sizeof(*this))
}

func (this *NodeResultMessage) Replicate() (ret MaybeRemoteMessage) {
	new := *this
	ret.Value(&new)
	return
}


