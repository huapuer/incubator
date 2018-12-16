package topo

import (
	"fmt"
	"github.com/incubator/common/maybe"
	"github.com/incubator/config"
	"github.com/incubator/interfaces"
	"net"
	"unsafe"
)

const (
	defaultSessionTopoClassName = "topo.defaultSessionTopo"
)

func init() {
	interfaces.RegisterTopoPrototype(defaultSessionTopoClassName, &defaultSessionTopo{}).Test()
}

type defaultSessionTopo struct {
	commonTopo

	hostSchema int32
	hosts      map[int64]interfaces.SessionHost
	hostCanon  interfaces.SessionHost
}

func (this defaultSessionTopo) New(attrs interface{}, cfg interfaces.Config) interfaces.IOC {
	ret := interfaces.MaybeTopo{}
	topo := &defaultSessionTopo{
		hosts: make(map[int64]interfaces.SessionHost),
	}

	topo.hostSchema = config.GetAttrInt32(attrs, "LocalHostSchema", config.CheckInt32GT0).Right()

	localHostCfg, ok := cfg.(*config.Config).HostMap[topo.hostSchema]
	if !ok {
		ret.Error(fmt.Errorf("no local host cfg found: %d", topo.hostSchema))
		return ret
	}
	topo.hostCanon = interfaces.GetHostPrototype(localHostCfg.Class).Right().(interfaces.SessionHost)

	ret.Value(topo)
	return ret
}

func (this *defaultSessionTopo) AddHost(id int64, conn net.Conn) (err maybe.MaybeError) {
	if _, ok := this.hosts[id]; ok {
		err.Error(fmt.Errorf("host already exsits: %d", id))
		return
	}
	host := this.hostCanon.Replicate().Right().(interfaces.SessionHost)
	host.SetPeer(conn)
	this.hosts[id] = host

	err.Error(nil)
	return
}

func (this defaultSessionTopo) SendToHost(id int64, msg interfaces.RemoteMessage) (err maybe.MaybeError) {
	host, ok := this.hosts[id]
	if !ok {
		err.Error(fmt.Errorf("no host found: %d", id))
		return
	}

	host.Receive(msg).Test()

	err.Error(nil)
	return
}

func (this defaultSessionTopo) LookupHost(id int64) (ret interfaces.MaybeHost) {
	panic("not implemented")
}

func (this defaultSessionTopo) SendToLink(hid int64, gid int64, msg interfaces.RemoteMessage) (err maybe.MaybeError) {
	panic("not implemented")
}

func (this defaultSessionTopo) TraverseOutLinksOfHost(hid int64, callback func(ptr unsafe.Pointer) bool) (err maybe.MaybeError) {
	panic("not implemented")
}

func (this defaultSessionTopo) GetRemoteHosts() []interfaces.Host {
	panic("not implemented")
}

func (this defaultSessionTopo) GetRemoteHostId(idx int32) int64 {
	panic("not implemented")
}

func (this defaultSessionTopo) Start() {}

func (this defaultSessionTopo) GetAddr() string {
	panic("not implemented")
}
