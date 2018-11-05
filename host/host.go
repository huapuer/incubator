package host

import (
	"errors"
	"fmt"
	"incubator/common/maybe"
	"incubator/config"
	"incubator/message"
)

var (
	hostsPrototype = make(map[string]Host)
)

func RegisterHostPrototype(name string, val Host) (err maybe.MaybeError) {
	if _, ok := hostsPrototype[name]; ok {
		err.Error(fmt.Errorf("host prototype redefined: %s", name))
		return
	}
	hostsPrototype[name] = val
	return
}

func GetHostPrototype(name string) (ret MaybeHost) {
	if prototype, ok := hostsPrototype[name]; ok {
		ret.Value(prototype)
		return
	}
	ret.Error(fmt.Errorf("host prototype for class not found: %s", name))
	return
}

type Host interface {
	config.IOC

	GetId() maybe.MaybeInt64
	SetId(int64) maybe.MaybeError
	Receive(message message.Message) maybe.MaybeError
}

type commonHost struct {
	id int64
}

func (this *commonHost) GetId() (id maybe.MaybeInt64) {
	if this.id < 0 {
		id.Error(errors.New("hostid less than 0."))
		return
	}
	id.Value(this.id)
	return
}

func (this *commonHost) SetId(id int64) (err maybe.MaybeError) {
	if id < 0 {
		err.Error(errors.New("hostid less than 0."))
		return
	}
	this.id = id
	return
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
