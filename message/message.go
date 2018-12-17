package message

import (
	"fmt"
	"github.com/incubator/common/maybe"
	"github.com/incubator/interfaces"
	"github.com/incubator/serialization"
)

////go:noescape
func RoutePackage(data []byte, layerId uint8, typ uint8) (err maybe.MaybeError) {
	l := interfaces.GetLayer(int32(layerId)).Right()
	msg := l.GetMessageCanonicalFromType(int32(typ)).Right()
	serialization.UnmarshalRemoteMessage(data, &msg).Test()
	router := l.GetRouter(int32(typ)).Right()
	router.Route(msg).Test()
	return
}

////go:noescape
func Route(m interfaces.RemoteMessage) (err maybe.MaybeError) {
	l := interfaces.GetLayer(int32(m.GetLayer())).Right()
	router := l.GetRouter(int32(m.GetType())).Right()
	router.Route(m).Test()

	err.Error(nil)
	return
}

func SendToHost(m interfaces.RemoteMessage) (err maybe.MaybeError) {
	if m.GetHostId() <= 0 {
		err.Error(fmt.Errorf("illegal host id: %d", m.GetHostId()))
	}
	interfaces.GetLayer(int32(m.GetLayer())).Right().GetTopo().SendToHost(m.GetHostId(), m).Test()
	return
}

func SendToLink(m interfaces.RemoteMessage, guestId int64) (err maybe.MaybeError) {
	if m.GetHostId() <= 0 {
		err.Error(fmt.Errorf("illegal host id: %d", m.GetHostId()))
	}
	if guestId <= 0 {
		err.Error(fmt.Errorf("illegal guest id: %d", guestId))
	}
	interfaces.GetLayer(int32(m.GetLayer())).Right().GetTopo().SendToLink(m.GetHostId(), guestId, m).Test()
	return
}

const (
	MASTER_NO = iota
	MASTER_YES
)

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
	interfaces.RemoteMessage

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
