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
)

var (
	messagePrototype = make(map[string]Message)
	messageCanonical = make(map[int]Message)
	messageRouters   = make(map[int]router.Router)
)

func RegisterMessagePrototype(name string, val Message) (err maybe.MaybeError){
	if _, ok := messagePrototype[name]; ok {
		err.Error(fmt.Errorf("message prototype redefined: %s", name))
		return
	}
	messagePrototype[name] = val
	return
}

func RegisterMessageCanonical(layerOffset int32, typ int32, cfg config.Config) (err maybe.MaybeError) {
	if typ <= 0 {
		err.Error(fmt.Errorf("illegal message type: %d", typ))
		return
	}
	if _, ok := messageCanonical[typ]; ok {
		err.Error(fmt.Errorf("message canonical type already exists: %d", typ))
		return
	}
	if _, ok := messageRouters[typ]; ok {
		err.Error(fmt.Errorf("router for message type already exists: %d", typ))
		return
	}
	if c, ok := cfg.Messages[typ]; ok {
		if p, ok := messagePrototype[c.Class]; ok {
			if rc, ok := cfg.Routers[c.RouterId]; ok {
				p.SetType(typ)
				messageCanonical[layerOffset + typ] = p
				r := router.GetRouter(layerOffset + rc.Id).Right()
				messageRouters[layerOffset + typ] = r
			}
			err.Error(fmt.Errorf("router config for message type not found: %d, %d", typ, c.RouterId))
			return
		}
		err.Error(fmt.Errorf("message prototype not found: %d", c.Class))
		return
	}
	err.Error(fmt.Errorf("message config not found: %d", typ))
	return
}

func GetMessageCanonical(typ int32) (ret MaybeMessage) {
	if msg, ok:=messageCanonical[typ];ok{
		ret.Value(msg)
		return
	}
	ret.Error(fmt.Errorf("message canonical does not exists: %d", typ))
	return
}

func RoutePackage(data []byte, typ int) (err maybe.MaybeError) {
	if m, ok := messageCanonical[typ]; ok {
		msg := m.Unmarshal(data).Right()
		msgType := msg.GetType().Right()
		if msgType != typ {
			err.Error(fmt.Errorf("message type unmatched with package: pkg:%d, msg:%d", typ, msgType))
			return
		}
		err = Route(msg)
		return
	}
	err.Error(fmt.Errorf("message canonical for type not found: %d", typ))
	return
}

func Route(m Message) (err maybe.MaybeError) {
	if r, ok := messageRouters[m.GetType().Right()]; ok {
		r.Route(m)
	}
	err.Error(fmt.Errorf("router for message type not found: %d", m.GetType().Right()))
	return
}

func SendTo(m Message, topoId int32, hostId int64) (err maybe.MaybeError) {
	if hostId < 0 {
		err.Error(fmt.Errorf("illegal host id: %d", hostId))
	}
	topo.GetTopo(topoId).Right().Lookup(hostId).Right().Receive(m).Test()
	return
}

type Message interface {
	SetType(int) maybe.MaybeError
	GetType() maybe.MaybeInt
	GetSize() int
	Process(context.Context) maybe.MaybeError
	GetHostId() maybe.MaybeInt64
	SetHostId(int64) maybe.MaybeError
	Marshal() []byte
	GetJsonBytes() maybe.MaybeBytes
	SetJsonField([]byte) maybe.MaybeError
	Unmarshal([]byte) MaybeMessage
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
	class.DefaultClass

	typ    int
	size   int
	hostId int64
}

func (this *commonMessage) Process(ctx context.Context) (err maybe.MaybeError) {
	err.Error(errors.New("calling abstract method:commonMessage.Process()"))
	return
}

func (this *commonMessage) GetType() (ret maybe.MaybeInt) {
	if this.typ <= 0 {
		ret.Error(fmt.Errorf("illegal message type: %d", this.typ))
		return
	}
	ret.Value(this.typ)
	return
}

func (this *commonMessage) SetType(typ int) (err maybe.MaybeError) {
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

func (this *commonMessage) SetHostId(seed int64) (err maybe.MaybeError) {
	if seed < 0 {
		err.Error(fmt.Errorf("illegal seed: %d", seed))
		return
	}
	this.hostId = seed
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

func (this *commonMessage) Marshal() (ret []byte) {
	derived := this.GetDerived()
	mi := (*mimicIFace)(unsafe.Pointer(&derived))

	msg := this.GetDerived().(Message)
	size := msg.GetSize()
	ms := &mimicSlice{mi.data, size, size}
	val := *(*[]byte)(unsafe.Pointer(ms))

	typ := msg.GetType().Right()
	jbytes := msg.GetJsonBytes().Right()
	jlth := len(jbytes)
	lth := len(val) + jlth + 2

	ret = append(ret, uint8(lth), uint8(len(val)), uint8(typ))
	ret = append(ret, val...)
	if jlth > 0 {
		ret = append(ret, jbytes...)
	}

	return
}

func (this *commonMessage) Unmarshal(data []byte) (msg MaybeMessage) {
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

	derived := this.GetDerived().(Message)

	mt := (*mimicIFace)(unsafe.Pointer(&derived))
	mi := (*mimicIFace)(unsafe.Pointer(&msg.value))
	mi.data = ms.addr
	mi.tab = mt.tab

	ljsn := lth - lval - 2
	if ljsn > 0 {
		jsn := data[lval+2 : lth]
		derived.SetJsonField(jsn).Test()
	}

	return
}
