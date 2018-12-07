package topo

import (
	"../common/maybe"
	"../config"
	"../host"
	"../message"
	"fmt"
	"net"
	"unsafe"
)

const (
	defaultSessionTopoClassName = "topo.defaultSessionTopo"
)

func init() {
	RegisterTopoPrototype(defaultSessionTopoClassName, &defaultSessionTopo{}).Test()
}

type defaultSessionTopo struct {
	hostSchema int32
	hosts      map[int64]host.SessionHost
	hostCanon  host.SessionHost
}

func (this defaultSessionTopo) New(attrs interface{}, cfg config.Config) config.IOC {
	ret := MaybeTopo{}
	topo := &defaultSessionTopo{
		hosts: make(map[int64]host.SessionHost),
	}

	topo.hostSchema = config.GetAttrInt32(attrs, "LocalHostSchema", config.CheckInt32GT0).Right()

	localHostCfg, ok := cfg.Hosts[topo.hostSchema]
	if !ok {
		ret.Error(fmt.Errorf("no local host cfg found: %d", topo.hostSchema))
		return ret
	}
	topo.hostCanon = host.GetHostPrototype(localHostCfg.Class).Right().(host.SessionHost)

	ret.Value(topo)
	return ret
}

func (this *defaultSessionTopo) AddHost(id int64, conn net.Conn) (err maybe.MaybeError) {
	if _, ok := this.hosts[id]; ok {
		err.Error(fmt.Errorf("host already exsits: %d", id))
		return
	}
	host := this.hostCanon.Replicate().Right().(host.SessionHost)
	host.SetPeer(conn)
	this.hosts[id] = host

	err.Error(nil)
	return
}

func (this defaultSessionTopo) SendToHost(id int64, msg message.RemoteMessage) (err maybe.MaybeError) {
	host, ok := this.hosts[id]
	if !ok {
		err.Error(fmt.Errorf("no host found: %d", id))
		return
	}

	host.Receive(msg).Test()

	err.Error(nil)
	return
}

func (this defaultSessionTopo) LookupHost(id int64) (ret host.MaybeHost) {
	panic("not implemented")
}

func (this defaultSessionTopo) SendToLink(hid int64, gid int64, msg message.RemoteMessage) (err maybe.MaybeError) {
	panic("not implemented")
}

func (this defaultSessionTopo) TraverseOutLinksOfHost(hid int64, callback func(ptr unsafe.Pointer) bool) (err maybe.MaybeError) {
	panic("not implemented")
}

func (this defaultSessionTopo) GetRemoteHosts() []host.Host {
	panic("not implemented")
}

func (this defaultSessionTopo) GetRemoteHostId(idx int32) int64 {
	panic("not implemented")
}

func (this defaultSessionTopo) Start() {
	panic("not implemented")
}

func (this defaultSessionTopo) GetLayer() int32 {
	panic("not implemented")
}

func (this *defaultSessionTopo) SetLayer(layer int32) {
	panic("not implemented")
}

func (this defaultSessionTopo) GetAddr() string {
	panic("not implemented")
}
