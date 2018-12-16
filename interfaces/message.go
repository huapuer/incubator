package interfaces

import (
	"fmt"
	"github.com/incubator/common/maybe"
)

var (
	messagePrototype = make(map[string]RemoteMessage)
)

func RegisterMessagePrototype(name string, val RemoteMessage) (err maybe.MaybeError) {
	if _, ok := messagePrototype[name]; ok {
		err.Error(fmt.Errorf("message prototype redefined: %s", name))
		return
	}
	messagePrototype[name] = val

	err.Error(nil)
	return
}

func GetMessagePrototype(name string) (ret MaybeRemoteMessage) {
	if msg, ok := messagePrototype[name]; ok {
		ret.Value(msg)
		return
	}
	ret.Error(fmt.Errorf("message prototype not found: %s", name))
	return
}

type Message interface {
	Process(Actor) maybe.MaybeError
}

type RemoteMessage interface {
	Message
	Serializable

	SetLayer(int8) maybe.MaybeError
	GetLayer() int8
	SetType(int8) maybe.MaybeError
	GetType() int8
	Master(int8)
	IsMaster() int8
	GetHostId() int64
	SetHostId(int64) maybe.MaybeError
}

type MaybeRemoteMessage struct {
	maybe.MaybeError

	value RemoteMessage
}

func (this *MaybeRemoteMessage) Value(value RemoteMessage) {
	this.Error(nil)
	this.value = value
}

func (this MaybeRemoteMessage) Right() RemoteMessage {
	this.Test()
	return this.value
}

type EchoMessage interface {
	Message

	SetSrcLayer(int8) maybe.MaybeError
	GetSrcLayer() int8
	GetSrcHostId() maybe.MaybeInt64
	SetSrcHostId(int64) maybe.MaybeError
}
