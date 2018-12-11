package host

import (
	"../common/maybe"
	"../config"
	"../message"
	"../storage"
	"fmt"
	"net"
)

var (
	hostPrototypes = make(map[string]Host)
)

func RegisterHostPrototype(name string, val Host) (err maybe.MaybeError) {
	if _, ok := hostPrototypes[name]; ok {
		err.Error(fmt.Errorf("host prototype redefined: %s", name))
		return
	}
	hostPrototypes[name] = val
	return
}

func GetHostPrototype(name string) (ret MaybeHost) {
	if prototype, ok := hostPrototypes[name]; ok {
		ret.Value(prototype)
		return
	}
	ret.Error(fmt.Errorf("host prototype for class not found: %s", name))
	return
}

type Host interface {
	config.IOC
	storage.DenseTableElement

	GetId() int64
	SetId(int64)
	Receive(message.RemoteMessage) maybe.MaybeError
	IsHealth() bool
	SetIP(string)
	SetPort(int)
	Start() maybe.MaybeError
}

type MaybeHost struct {
	config.IOC

	maybe.MaybeError
	value Host
}

func (this MaybeHost) Value(value Host) {
	this.Error(nil)
	this.value = value
}

func (this MaybeHost) Right() Host {
	this.Test()
	return this.value
}

func (this MaybeHost) New(cfg config.Config, args ...int32) config.IOC {
	panic("not implemented.")
}

type commonHost struct {
	id int64
}

func (this commonHost) GetId() int64 {
	return this.id
}

func (this *commonHost) SetId(id int64) {
	this.id = id
	return
}

func (this commonHost) SetIP(string) {
	panic("not implemented")
}

func (this commonHost) SetPort(int) {
	panic("not implemented")
}

func (this commonHost) Start() maybe.MaybeError {
	panic("not implemented")
}

type commonLinkHost struct {
	guestId int64
}

func (this *commonLinkHost) GetGuestId() int64 {
	return this.guestId
}

func (this *commonLinkHost) SetGuestId(id int64) {
	this.guestId = id
	return
}

type SessionHost interface {
	Host

	SetPeer(net.Conn)
	Replicate() MaybeHost
}
