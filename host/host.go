package host

import (
	"../common/maybe"
	"../config"
	"../message"
	"fmt"
	"../storage"
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

	GetId() int64
	SetId(int64)
	Valid(bool)
	IsValid() bool
	Receive(net.Conn, message.RemoteMessage) maybe.MaybeError
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

type LocalHost interface {
	Host

	storage.DenseTableElement
}

type commonHost struct {
	id    int64
	valid bool
}

func (this *commonHost) GetId() int64 {
	return this.id
}

func (this *commonHost) SetId(id int64) {
	this.id = id
	return
}

func (this *commonHost) IsValid() bool {
	return this.valid
}

func (this *commonHost) Valid(v bool) {
	this.valid = v
	return
}
