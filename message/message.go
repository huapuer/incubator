package message

import (
	"../actor"
	"../common/maybe"
	"../layer"
	"../serialization"
	"fmt"
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

//go:noescape
func RoutePackage(data []byte, layerId uint8, typ uint8) (err maybe.MaybeError) {
	l := layer.GetLayer(int32(layerId)).Right()
	msg := l.GetMessageCanonicalFromType(int32(typ)).Right()
	serialization.UnmarshalRemoteMessage(data, msg).Test()
	router := l.GetRouter(int32(typ)).Right()
	router.Route(msg).Test()
	return
}

//go:noescape
func Route(m RemoteMessage) (err maybe.MaybeError) {
	l := layer.GetLayer(int32(m.GetLayer())).Right()
	router := l.GetRouter(int32(m.GetType())).Right()
	router.Route(m).Test()
	return
}

func SendToHost(m RemoteMessage, layerId int32, hostId int64) (err maybe.MaybeError) {
	if hostId <= 0 {
		err.Error(fmt.Errorf("illegal host id: %d", hostId))
	}
	layer.GetLayer(layerId).Right().LookupHost(hostId).Right().Receive(m).Test()
	return
}

func SendToLink(m RemoteMessage, layerId int32, hostId int64, guestId int64) (err maybe.MaybeError) {
	if hostId <= 0 {
		err.Error(fmt.Errorf("illegal host id: %d", hostId))
	}
	if guestId <= 0 {
		err.Error(fmt.Errorf("illegal guest id: %d", guestId))
	}
	layer.GetLayer(layerId).Right().LookupLink(hostId, guestId).Right().Receive(m).Test()
	return
}

type Message interface {
	Process(actor.Actor) maybe.MaybeError
}

type RemoteMessage interface {
	Message
	serialization.Serializable

	SetLayer(int8) maybe.MaybeError
	GetLayer() int8
	SetType(int8) maybe.MaybeError
	GetType() int8
	Master(int8)
	IsMaster() int8
	GetHostId() int64
	SetHostId(int64) maybe.MaybeError
	Replicate() MaybeRemoteMessage
}

type MaybeRemoteMessage struct {
	maybe.MaybeError

	value RemoteMessage
}

func (this MaybeRemoteMessage) Value(value RemoteMessage) {
	this.Error(nil)
	this.value = value
}

func (this MaybeRemoteMessage) Right() RemoteMessage {
	this.Test()
	return this.value
}

type commonMessage struct {
	layerId int8
	typ     int8
	master  int8
	hostId  int64
}

func (this *commonMessage) GetLayer() int8 {
	return this.layerId
}

func (this *commonMessage) SetLayer(layer int8) (err maybe.MaybeError) {
	if layer < 0 {
		err.Error(fmt.Errorf("illegal message layerId: %d", layer))
		return
	}
	this.layerId = layer
	return
}

func (this *commonMessage) GetType() int8 {
	return this.typ
}

func (this *commonMessage) SetType(typ int8) (err maybe.MaybeError) {
	if typ <= 0 {
		err.Error(fmt.Errorf("illegal message type: %d", typ))
		return
	}
	this.typ = typ
	return
}

func (this *commonMessage) IsMaster() int8 {
	return this.master
}

func (this *commonMessage) Master(b int8) {
	this.master = b
	return
}

func (this *commonMessage) GetHostId() int64 {
	return this.hostId
}

func (this *commonMessage) SetHostId(hostId int64) (err maybe.MaybeError) {
	if hostId < 0 {
		err.Error(fmt.Errorf("illegal seed: %d", hostId))
		return
	}
	this.hostId = hostId
	return
}

type EchoMessage interface {
	Message

	SetSrcLayer(int8) maybe.MaybeError
	GetSrcLayer() int8
	GetSrcHostId() maybe.MaybeInt64
	SetSrcHostId(int64) maybe.MaybeError
}

type commonEchoMessage struct {
	commonMessage

	srcLayer  int8
	srcHostId int64
}

func (this *commonEchoMessage) GetSrcLayer() int8 {
	return this.srcLayer
}

func (this *commonEchoMessage) SetSrcLayer(layer int8) (err maybe.MaybeError) {
	if layer < 0 {
		err.Error(fmt.Errorf("illegal message layerId: %d", layer))
		return
	}
	this.srcLayer = layer
	return
}

func (this *commonEchoMessage) GetSrcHostId() int64 {
	return this.srcHostId
}

func (this *commonEchoMessage) SetSrcHostId(hostId int64) (err maybe.MaybeError) {
	if hostId < 0 {
		err.Error(fmt.Errorf("illegal seed: %d", hostId))
		return
	}
	this.srcHostId = hostId
	return
}

type LinkMessage interface {
	RemoteMessage

	SetGuestId(int64)
	GetGuestId() int64
}

type commonLinkMessage struct {
	guestId int64
}

func (this *commonLinkMessage) GetGuestId() int64 {
	return this.guestId
}

func (this *commonLinkMessage) SetGuestId(id int64) {
	this.guestId = id
	return
}
