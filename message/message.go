package message

import (
	"../actor"
	"../common/maybe"
	"../serialization"
	"../topo"
	"context"
	"errors"
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

func RoutePackage(data []byte, layer uint8, typ uint8) (err maybe.MaybeError) {
	tp := topo.GetTopo(int32(layer)).Right()
	msgCanon := tp.GetMessageCanonicalFromType(int32(typ)).Right()
	msg := msgCanon.Duplicate().Right()
	serialization.Unmarshal(data, msg).Test()
	router := tp.GetRouter(int32(typ)).Right()
	router.Route(msg).Test()
	return
}

func Route(m RemoteMessage) (err maybe.MaybeError) {
	tp := topo.GetTopo(int32(m.GetLayer())).Right()
	router := tp.GetRouter(int32(m.GetType())).Right()
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
	Duplicate() MaybeRemoteMessage
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
	layer  int8
	typ    int8
	master int8
	hostId int64
}

func (this *commonMessage) Process(ctx context.Context) (err maybe.MaybeError) {
	err.Error(errors.New("calling abstract method:commonMessage.Process()"))
	return
}

func (this *commonMessage) GetLayer() int8 {
	return this.layer
}

func (this *commonMessage) SetLayer(layer int8) (err maybe.MaybeError) {
	if layer < 0 {
		err.Error(fmt.Errorf("illegal message layer: %d", layer))
		return
	}
	this.layer = layer
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

func (this commonMessage) copyPaste(msg RemoteMessage) {
	msg.SetType(this.typ)
	msg.SetLayer(this.layer)
}

type SeesionedMessage interface {
	RemoteMessage

	SetSessionId(int64)
	GetSesseionId() int64
	ToServer()
	ToClient()
	IsToServer() maybe.MaybeBool
}

type commonSessionedMessage struct {
	commonMessage

	toServer  int8
	sessionId int64
}

const (
	SESSEION_MESSAGE_TO_SERVER = iota
	SESSION_MESSAGE_TO_CLIENT
)

func (this commonSessionedMessage) GetSessionId() int64 {
	return this.sessionId
}

func (this *commonSessionedMessage) SetSessionId(sid int64) {
	this.sessionId = sid
}

func (this *commonSessionedMessage) ToServer() {
	this.toServer = SESSEION_MESSAGE_TO_SERVER
}

func (this *commonSessionedMessage) ToClient() {
	this.toServer = SESSION_MESSAGE_TO_CLIENT
}

func (this commonSessionedMessage) IsToServer() (ret maybe.MaybeBool) {
	if this.toServer < 0 {
		ret.Error(fmt.Errorf("is to server flag not set"))
		return
	}
	ret.Value(this.toServer == SESSEION_MESSAGE_TO_SERVER)
	return
}

type EchoMessage interface {
	SeesionedMessage

	SetSrcLayer(int8) maybe.MaybeError
	GetSrcLayer() int8
	GetSrcHostId() maybe.MaybeInt64
	SetSrcHostId(int64) maybe.MaybeError
}

type commonEchoMessage struct {
	commonSessionedMessage

	srcLayer  int8
	srcHostId int64
}

func (this *commonEchoMessage) GetSrcLayer() int8 {
	return this.srcLayer
}

func (this *commonEchoMessage) SetSrcLayer(layer int8) (err maybe.MaybeError) {
	if layer < 0 {
		err.Error(fmt.Errorf("illegal message layer: %d", layer))
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
