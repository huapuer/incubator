package host

import (
	"../common/maybe"
	"../config"
	"../message"
	"../network"
	"errors"
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

	address string
	client  network.Client
}

func (this *defaultRemoteHost) Receive(msg message.Message) (err maybe.MaybeError) {
	this.client.Send(msg).Test()
	return
}

func (this defaultRemoteHost) New(attrs interface{}, cfg config.Config) config.IOC {
	ret := MaybeHost{}

	if attrs == nil {
		ret.Error(fmt.Errorf("attrs is nil when new host: %s", defaultRemoteHostClassName))
		return ret
	}
	attrsMap, ok := attrs.(map[string]interface{})
	if !ok {
		ret.Error(fmt.Errorf("illegal cfg type when new host: %s", defaultRemoteHostClassName))
		return ret
	}

	clientSchema, ok := attrsMap["ClientSchema"]
	if !ok {
		ret.Error(errors.New("attribute ClientSchema not found"))
		return ret
	}
	clientSchemaInt, ok := clientSchema.(int32)
	if !ok {
		ret.Error(fmt.Errorf("client schema cfg type error(expecting int): %+v", clientSchema))
		return ret
	}

	clientCfg, ok := cfg.Clients[clientSchemaInt]
	if !ok {
		ret.Error(fmt.Errorf("client cfg not found: %d", clientCfg))
		return ret
	}

	//TODO: real logic
	ret.Value(&defaultRemoteHost{
		commonHost: commonHost{
			valid: true,
		},
		client: network.DefaultClient.New(clientCfg.Attributes, cfg).(network.MaybeDefualtClient).Right(),
	})
	return ret
	return ret
}

func (this defaultRemoteHost) GetJsonBytes() (ret maybe.MaybeBytes) {
	ret.Value([]byte{})
	return
}

func (this *defaultRemoteHost) SetJsonField(data []byte) (err maybe.MaybeError) {
	err.Error(nil)
	return
}

func (this defaultRemoteHost) GetSize() int32 {
	return int32(unsafe.Sizeof(this))
}
