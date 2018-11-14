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
	msg := msgCanon.Unmarshal(data, msgCanon).Right()
	router := tp.GetRouter(typ).Right()
	router.Route(msg).Test()
	return
}

func Route(m Message) (err maybe.MaybeError) {
	tp:=topo.GetTopo(m.GetType()).Right()
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
	SetLayer(uint8) maybe.MaybeError
	GetLayer() uint8
	SetType(uint8) maybe.MaybeError
	GetType() uint8
	GetSize() int32
	Process(context.Context) maybe.MaybeError
	GetHostId() maybe.MaybeInt64
	SetHostId(int64) maybe.MaybeError
	Marshal(Message) []byte
	GetJsonBytes() maybe.MaybeBytes
	SetJsonField([]byte) maybe.MaybeError
	Unmarshal([]byte, Message) MaybeMessage
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

type mimicSlice struct {
	addr *unsafe.ArbitraryType
	len  int
	cap  int
}

type mimicIFace struct {
	tab  unsafe.Pointer
	data unsafe.Pointer
}

func (this *commonMessage) Marshal(msg Message) (ret []byte) {
	mi := (*mimicIFace)(unsafe.Pointer(&msg))

	size := msg.GetSize()
	ms := &mimicSlice{mi.data, size, size}
	val := *(*[]byte)(unsafe.Pointer(ms))

	jbytes := msg.GetJsonBytes().Right()

	lth := int32(len(val) + len(jbytes) + 1 + unsafe.Sizeof(int32(0)))

	ret = append(ret, uint8(lth), uint8(len(val)))
	ret = append(ret, val...)
	if len(jbytes) > 0 {
		ret = append(ret, jbytes...)
	}

	return
}

func (this *commonMessage) Unmarshal(data []byte, canon Message) (msg MaybeMessage) {
	l := len(data)
	if l < 4 {
		msg.Error(fmt.Errorf("message bytes too short: %d", l))
		return
	}
	lth := int(data[0])
	if lth < 0 {
		msg.Error(fmt.Errorf("message claims negative lenth: %d", lth))
		return
	}
	if l != lth {
		msg.Error(fmt.Errorf("message length not equal to claimed, %d != %d", l, lth))
		return
	}

	lval := int(data[1])
	if lth < 0 {
		msg.Error(fmt.Errorf("message claims negative binary length: %d", lval))
		return
	}
	if lth < lval+3 {
		msg.Error(fmt.Errorf("message length shorter than claimed binary length + header length, %d  < %d + 3", l, lval))
		return
	}

	val := data[2 : lval+2]

	ms := (*mimicSlice)(unsafe.Pointer(&val))

	mt := (*mimicIFace)(unsafe.Pointer(&canon))
	mi := (*mimicIFace)(unsafe.Pointer(&msg.value))
	mi.data = ms.addr
	mi.tab = mt.tab

	ljsn := lth - lval - 2
	if ljsn > 0 {
		jsn := data[lval+2 : lth]
		canon.SetJsonField(jsn).Test()
	}

	return
}

func (this commonMessage) copyPaste (msg Message) {
	msg.SetType(this.typ)
	msg.SetLayer(this.layer)
}
