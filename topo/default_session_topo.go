package topo

import (
	"../common/maybe"
	"../config"
	"../host"
	"errors"
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

	attrsMap, ok := attrs.(map[string]interface{})
	if !ok {
		ret.Error(fmt.Errorf("illegal cfg type when new layer %s", defaultTopoClassName))
		return ret
	}

	localHostSchema, ok := attrsMap["LocalHostSchema"]
	if !ok {
		ret.Error(errors.New("attribute LocalHostSchema not found"))
		return ret
	}
	localHostSchemaInt, ok := localHostSchema.(int32)
	if !ok {
		ret.Error(fmt.Errorf("local host class cfg type error(expecting int32): %+v", localHostSchema))
		return ret
	}
	if localHostSchemaInt <= 0 {
		ret.Error(fmt.Errorf("illegal LocalHostSchema: %d", localHostSchemaInt))
		return ret
	}
	topo.hostSchema = localHostSchemaInt

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

func (this defaultSessionTopo) LookupHost(id int64) (ret host.MaybeHost) {
	host, ok := this.hosts[id]
	if !ok {
		ret.Error(fmt.Errorf("no host found: %d", id))
		return
	}
	ret.Value(host)
	return
}

func (this defaultSessionTopo) LookupLink(hid int64, gid int64) (ret host.MaybeHost) {
	panic("not implemented")
}

func (this defaultSessionTopo) TraverseLinksOfHost(hid int64, callback func(ptr unsafe.Pointer) bool) (err maybe.MaybeError) {
	panic("not implemented")
}

func (this defaultSessionTopo) GetRemoteHosts() []host.Host {
	panic("not implemented")
}
