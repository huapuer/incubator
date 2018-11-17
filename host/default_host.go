package host

import (
	"../common/maybe"
	"../config"
	"../message"
	"../network"
	"errors"
	"fmt"
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
	message.Route(msg).Test()
	return
}

func (this defaultLocalHost) New(attrs interface{}, cfg config.Config) config.IOC {
	ret := MaybeHost{}
	//TODO: real logic
	ret.Value(&defaultLocalHost{
		commonHost{
			valid:true,
		},
	})
	return ret
}

type defaultRemoteHost struct {
	commonHost

	address string
	client  network.Client
}

func (this *defaultRemoteHost) Receive(msg message.Message) (err maybe.MaybeError) {
	msg.SetHostId(this.GetId().Right()).Test()
	this.client.Send(msg).Test()
	return
}

func (this defaultRemoteHost) New(attrs interface{}, cfg config.Config) config.IOC {
	ret := MaybeHost{}

	if attrs == nil {
		ret.Error(fmt.Errorf("attrs is nil when new host: %s", defaultRemoteHostClassName))
		return
	}
	attrsMap, ok := attrs.(map[string]interface{})
	if !ok {
		ret.Error(fmt.Errorf("illegal cfg type when new host: %s", defaultRemoteHostClassName))
		return
	}

	clientSchema, ok := attrsMap["ClientSchema"]
	if !ok {
		ret.Error(errors.New("attribute ClientSchema not found"))
		return
	}
	clientSchemaInt, ok := clientSchema.(int32)
	if !ok {
		ret.Error(fmt.Errorf("client schema cfg type error(expecting int): %+v", clientSchema))
		return
	}

	clientCfg, ok := cfg.Clients[clientSchemaInt]
	if !ok {
		ret.Error(fmt.Errorf("client cfg not found: %d", clientCfg))
		return
	}

	//TODO: real logic
	ret.Value(&defaultRemoteHost{
		commonHost{
			valid:true,
		},
		client: network.DefaultClient.New(clientCfg.Attributes, cfg).(network.MaybeDefualtClient).Right(),
	})
	return ret
	return ret
}
