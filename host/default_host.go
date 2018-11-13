package host

import (
	"../common/maybe"
	"../config"
	"../message"
	"../network"
)

const (
	defaultLocalHostClassName  = "actor.defaultLocalHost"
	defaultRemoteHostClassName = "actor.defaultRemoteHost"
)

func init() {
	RegisterHostPrototype(defaultLocalHostClassName, &defaultLocalHost{}).Test()
	RegisterHostPrototype(defaultRemoteHostClassName, &defaultRemoteHost{}).Test()
}

type defaultLocalHost struct {
	commonHost
}

func (this *defaultLocalHost) Receive(msg message.Message) (err maybe.MaybeError) {
	msg.SetHostId(this.GetId().Right()).Test()
	err = message.Route(msg)
	return
}

func (this defaultLocalHost) New(attrs interface{}, cfg config.Config) config.IOC {
	ret := MaybeHost{}
	//TODO: real logic
	ret.Value(&defaultLocalHost{})
	return ret
}

type defaultRemoteHost struct {
	commonHost

	address string
	client  network.Client
}

func (this *defaultRemoteHost) Receive(msg message.Message) (err maybe.MaybeError) {
	err = msg.SetHostId(this.id)
	this.client.Send(msg).Test()
	return
}

func (this defaultRemoteHost) New(attrs interface{}, cfg config.Config) config.IOC {
	ret := MaybeHost{}
	//TODO: real logic
	ret.Value(&defaultLocalHost{})
	return ret
	return ret
}
