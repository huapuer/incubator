package host

import (
	"fmt"
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

func (this defaultLocalHost) New(cfg config.Config) config.IOC {
	ret := MaybeHost{}
	if host, ok := cfg.Hosts[defaultLocalHostClassName]; ok {
		//TODO: real logic
		host = host
		ret.Value(&defaultLocalHost{})
		return ret
	}
	ret.Error(fmt.Errorf("no host class cfg found: %s", defaultLocalHostClassName))
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

func (this defaultRemoteHost) New(cfg config.Config) config.IOC {
	ret := MaybeHost{}
	if host, ok := cfg.Hosts[defaultRemoteHostClassName]; ok {
		//TODO: real logic
		host = host
		ret.Value(&defaultLocalHost{})
		return ret
	}
	ret.Error(fmt.Errorf("no host class cfg found: %s", defaultRemoteHostClassName))
	return ret
}
