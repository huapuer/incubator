package host

import (
	"errors"
	"fmt"
	"github.com/incubator/common/maybe"
	"github.com/incubator/config"
	"github.com/incubator/interfaces"
	"github.com/incubator/network"
	"time"
	"unsafe"
)

const (
	defaultRemoteHostClassName = "host.defaultRemoteHost"
)

func init() {
	RegisterHostPrototype(defaultRemoteHostClassName, &defaultRemoteHost{}).Test()
}

type defaultRemoteHost struct {
	commonHost
	defaultHealthManager

	ip     string
	port   int
	client interfaces.Client
}

func (this *defaultRemoteHost) Receive(msg interfaces.RemoteMessage) (err maybe.MaybeError) {
	this.client.Send(msg).Test()
	return
}

func (this defaultRemoteHost) New(attrs interface{}, cfg interfaces.Config) interfaces.IOC {
	ret := interfaces.MaybeHost{}

	clientSchema := config.GetAttrInt32(attrs, "ClientSchema", nil).Right()
	checkIntvl := config.GetAttrInt64(attrs, "CheckIntvl", config.CheckInt64GT0).Right()
	heartbeatIntvl := config.GetAttrInt64(attrs, "HeartbeatIntvl", config.CheckInt64GT0).Right()

	clientCfg, ok := cfg.(*config.Config).ClientMap[clientSchema]
	if !ok {
		ret.Error(fmt.Errorf("client cfg not found: %d", clientSchema))
		return ret
	}

	host := &defaultRemoteHost{
		client: network.GetClientPrototype(clientCfg.Class).
			Right().New(clientCfg.Attributes, cfg).(interfaces.MaybeClient).Right(),
		defaultHealthManager: defaultHealthManager{
			health:         true,
			checkIntvl:     time.Duration(checkIntvl),
			heartbeatIntvl: time.Duration(heartbeatIntvl),
		},
	}

	ret.Value(host)
	return ret
}

func (this defaultRemoteHost) IsHealth() bool {
	return this.health
}

func (this defaultRemoteHost) SetIP(ip string) {
	this.ip = ip
}

func (this defaultRemoteHost) SetPort(port int) {
	this.port = port
}

func (this defaultRemoteHost) Start() (err maybe.MaybeError) {
	if this.ip == "" {
		err.Error(errors.New("remote host ip not set"))
		return
	}
	if this.port <= 0 {
		err.Error(fmt.Errorf("illegal remote host port: %d", this.port))
		return
	}
	if this.client == nil {
		err.Error(errors.New("remote host client not inited"))
		return
	}

	this.client.Connect(fmt.Sprint("%s:%d", this.ip, this.port))

	err.Error(nil)
	return
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
