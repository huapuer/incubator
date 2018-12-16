package interfaces

import (
	"fmt"
	"github.com/incubator/common/maybe"
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

	err.Error(nil)
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
	IOC
	DenseTableElement

	GetId() int64
	SetId(int64)
	Receive(RemoteMessage) maybe.MaybeError
	IsHealth() bool
	SetIP(string)
	SetPort(int)
	Start() maybe.MaybeError
}

type MaybeHost struct {
	IOC

	maybe.MaybeError
	value Host
}

func (this *MaybeHost) Value(value Host) {
	this.Error(nil)
	this.value = value
}

func (this MaybeHost) Right() Host {
	this.Test()
	return this.value
}

func (this MaybeHost) New(attr interface{}, cfg Config) IOC {
	panic("not implemented.")
}

type SessionHost interface {
	Host

	SetPeer(net.Conn)
	Replicate() MaybeHost
}
