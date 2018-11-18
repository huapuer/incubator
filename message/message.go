package message

import (
	"context"
	"errors"
	"fmt"
	"../common/class"
	"../common/maybe"
	"../config"
	"../router"
	"unsafe"
	"../topo"
	"strings"
	"github.com/incubator/serialization"
	"github.com/incubator/actor"
	"github.com/incubator/host"
)

var (
	messagePrototype = make(map[string]Message)
)

func RegisterMessagePrototype(name string, val Message) (err maybe.MaybeError){
	if _, ok := messagePrototype[name]; ok {
		err.Error(fmt.Errorf("message prototype redefined: %s", name))
		return
	}
	messagePrototype[name] = val
	return
}

func GetMessagePrototype(name string) (ret MaybeMessage) {
	if msg, ok := messagePrototype[name]; ok{
		ret.Value(msg)
		return
	}
	ret.Error(fmt.Errorf("message prototype not found: %s", name))
	return
}

func RoutePackage(data []byte, layer uint8, typ uint8) (err maybe.MaybeError) {
	tp := topo.GetTopo(layer).Right()
	msgCanon := tp.GetMessageCanonicalFromType(typ).Right()
	msg := msgCanon.Duplicate().Right()
	serialization.Unmarshal(data, msg).Test()
	router := tp.GetRouter(typ).Right()
	router.Route(msg).Test()
	return
}

func Route(m Message) (err maybe.MaybeError) {
	tp:=topo.GetTopo(m.GetLayer()).Right()
	router := tp.GetRouter(m.GetType()).Right()
	router.Route(m).Test()
	return
}

func SendTo(m Message, topoId int32, hostId int64) (err maybe.MaybeError) {
	if hostId <= 0 {
		err.Error(fmt.Errorf("illegal host id: %d", hostId))
	}
	topo.GetTopo(topoId).Right().Lookup(hostId).Right().Receive(m).Test()
	return
}

type Message interface {
	serialization.Serializable

	SetLayer(uint8) maybe.MaybeError
	GetLayer() uint8
	SetType(uint8) maybe.MaybeError
	GetType() uint8
	Master(bool)
	IsMaster() bool
	Process(actor.Actor) maybe.MaybeError
	GetHostId() maybe.MaybeInt64
	SetHostId(int64) maybe.MaybeError
	Duplicate() MaybeMessage
}

type MaybeMessage struct {
	maybe.MaybeError

	value Message
}

func (this MaybeMessage) Value(value Message) {
	this.Error(nil)
	this.value = value
}

func (this MaybeMessage) Right() Message {
	this.Test()
	return this.value
}

type commonMessage struct {
	layer  uint8
	typ    uint8
	master bool
	hostId int64
}

func (this *commonMessage) Process(ctx context.Context) (err maybe.MaybeError) {
	err.Error(errors.New("calling abstract method:commonMessage.Process()"))
	return
}

func (this *commonMessage) GetLayer() int32 {
	return this.layer
}

func (this *commonMessage) SetLayer(layer uint8) (err maybe.MaybeError) {
	if layer <= 0 {
		err.Error(fmt.Errorf("illegal message layer: %d", layer))
		return
	}
	this.layer = layer
	return
}

func (this *commonMessage) GetType() int32 {
	return this.typ
}

func (this *commonMessage) SetType(typ uint8) (err maybe.MaybeError) {
	if typ <= 0 {
		err.Error(fmt.Errorf("illegal message type: %d", typ))
		return
	}
	this.typ = typ
	return
}

func (this *commonMessage) IsMater() bool {
	return this.master
}

func (this *commonMessage) Master(b bool) {
	this.master = b
	return
}

func (this *commonMessage) GetHostId() (ret maybe.MaybeInt64) {
	if this.hostId < 0 {
		ret.Error(errors.New("hostid less than 0."))
		return
	}
	ret.Value(this.hostId)
	return
}

func (this *commonMessage) SetHostId(hostId int64) (err maybe.MaybeError) {
	if hostId < 0 {
		err.Error(fmt.Errorf("illegal seed: %d", hostId))
		return
	}
	this.hostId = hostId
	return
}

func (this commonMessage) copyPaste (msg Message) {
	msg.SetType(this.typ)
	msg.SetLayer(this.layer)
}
