package host

import (
	"../common/maybe"
	"../config"
	"../message"
	"../network"
	"fmt"
	"unsafe"
)

const (
	defaultRemoteHostClassName = "actor.defaultRemoteHost"
)

func init() {
	RegisterHostPrototype(defaultRemoteHostClassName, &defaultRemoteHost{}).Test()
}

type defaultRemoteHost struct {
	commonHost

	client  network.Client
}

func (this *defaultRemoteHost) Receive(msg message.RemoteMessage) (err maybe.MaybeError) {
	this.client.Send(msg).Test()
	return
}

func (this defaultRemoteHost) New(attrs interface{}, cfg config.Config) config.IOC {
	ret := MaybeHost{}

	clientSchema := config.GetAttrInt32(attrs, "ClientSchema", nil).Right()
	addr := config.GetAttrString(attrs, "Address", config.CheckStringNotEmpty).Right()

	clientCfg, ok := cfg.Clients[clientSchema]
	if !ok {
		ret.Error(fmt.Errorf("client cfg not found: %d", clientSchema))
		return ret
	}

	host := &defaultRemoteHost{
		commonHost: commonHost{
			valid: true,
		},
		client: network.DefaultClient.New(clientCfg.Attributes, cfg).(network.MaybeDefualtClient).Right(),
	}
	host.client.Connect(addr)

	ret.Value(host)
	return ret
}

func (this defaultRemoteHost) GetSize() int32 {
	panic("not implemented")
}

func (this defaultRemoteHost) Get(key int64, ptr unsafe.Pointer) bool {
	panic("not implemented")
}

func (this defaultRemoteHost) Put(dst unsafe.Pointer, src unsafe.Pointer) bool {
	panic("not implemented")
}

func (this defaultRemoteHost) Erase(key int64, ptr unsafe.Pointer) bool {
	panic("not implemented")
}
