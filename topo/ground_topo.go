package topo

import (
	"../config"
	"../host"
	"errors"
	"fmt"
	"../common/maybe"
	"unsafe"
)

const (
	groundTopoClassName = "topo.groundTopo"
)

func init() {
	RegisterTopoPrototype(groundTopoClassName, &groundTopo{}).Test()
}

type groundTopo struct {
	hostSchema int32
	host host.Host
}

func (this groundTopo) New(attrs interface{}, cfg config.Config) config.IOC {
	ret := MaybeTopo{}
	topo := &groundTopo{}

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
	topo.host = host.GetHostPrototype(localHostCfg.Class).Right().(host.LocalHost)

	ret.Value(topo)
	return ret
}

func (this groundTopo) LookupHost(id int64) (ret host.MaybeHost) {
	if this.host == nil {
		ret.Error(errors.New("no host found"))
		return
	}
	ret.Value(this.host)
	return
}

func (this groundTopo) LookupLink(hid int64, gid int64) (ret host.MaybeHost) {
	panic("not implemented")
}

func (this groundTopo) TraverseLinksOfHost(hid int64, callback func(ptr unsafe.Pointer) bool) (err maybe.MaybeError) {
	panic("not implemented")
}

func (this groundTopo) GetRemoteHosts() []host.Host {
	panic("not implemented")
}
