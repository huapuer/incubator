package message

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"incubator/common/class"
	"incubator/common/maybe"
	"incubator/config"
	"incubator/router"
	"os"
	"unsafe"
)

var (
	messagePrototype = make(map[string]Message)
	messageCanonical = make(map[int]Message)
	messageRouters   = make(map[int]router.Router)
)

func RegisterMessagePrototype(name string, val Message) {
	if _, ok := messagePrototype[name]; ok {
		logrus.Errorf("message prototype redefined: %s", name)
		os.Exit(1)
	}
	messagePrototype[name] = val
}

func RegisterMessageCanonical(cfg config.Config, typ int) (err maybe.MaybeError) {
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
		if p, ok := messagePrototype[c.ClassName]; ok {
			if rc, ok := cfg.Routers[c.Router]; ok {
				p.SetType(typ)
				messageCanonical[typ] = p
				r := router.GetRouter(rc.ClassName).Right()
				messageRouters[typ] = r
			}
			err.Error(fmt.Errorf("router config for message type not found: %d, %d", typ, c.Router))
			return
		}
		err.Error(fmt.Errorf("message prototype not found: %d", c.ClassName))
		return
	}
	err.Error(fmt.Errorf("message config not found: %d", typ))
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

type Message interface {
	SetType(int) maybe.MaybeError
	GetType() maybe.MaybeInt
	GetSize() maybe.MaybeInt
	Process(context.Context) maybe.MaybeError
	GetHostId() maybe.MaybeInt64
	SetHostId(int64) maybe.MaybeError
	Marshal() []byte
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
	Addr *unsafe.ArbitraryType
	len  int
	cap  int
}

func (this *commonMessage) Marshal() (ret []byte) {
	size := this.GetDerived().(Message).GetSize().Right()
	ms := &mimicSlice{unsafe.Pointer(this), size, size}
	val := *(*[]byte)(unsafe.Pointer(ms))
	typ := this.GetDerived().(Message).GetType().Right()
	lth := len(val) + 2
	ret = append(ret, uint8(lth), uint8(typ))
	ret = append(ret, val...)
	return
}

func (this *commonMessage) Unmarshal(ret []byte) (msg MaybeMessage) {
	msg.Error(errors.New("calling abstract method:commonMessage.Unmarshal()"))
	return
}
