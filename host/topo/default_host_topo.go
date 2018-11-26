package topo

import (
	"incubator/storage"
	"incubator/config"
	"incubator/host"
	"fmt"
	"incubator/serialization"
)

const (
	defaultHostTopoClassName = "actor.defaultLocalHost"
)

func init() {
	RegisterHostTopoPrototype(defaultHostTopoClassName, &defaultHostTopo{}).Test()
}

type defaultHostTopo struct {
	localHosts       storage.DenseTable
}

func (this defaultHostTopo) New(attrs interface{}, cfg config.Config) config.IOC {
	ret := MaybeHostTopo{}



	return ret
}

func (this defaultHostTopo) Lookup(id int64) (ret host.MaybeHost) {
	mod := int32(id % (int64(this.remoteNum)))
	idx := int32(id/int64(this.remoteNum)/int64(this.backupFactor+1)) + mod
	hosts := make([]host.Host, 0, 0)

	if mod == this.localHostMod {
		if idx > int32(this.localHosts.ElemLen()) {
			ret.Error(fmt.Errorf("master id exceeds local host range: %d", id))
			return
		}
		h := this.localHostCanon
		ptr := this.localHosts.Aquire(0, int64(idx)).Right()
		serialization.Ptr2IFace(&h, ptr)
		hosts = append(hosts, h)
	} else {
		hosts = append(hosts, this.remoteHosts[mod])
	}
	if mod < this.localHostMod+this.backupFactor {
		if idx > int32(this.localHosts.ElemLen()) {
			ret.Error(fmt.Errorf("slave id exceeds local host range: %d", id))
			return
		}
		h := this.localHostCanon
		ptr := this.localHosts.Aquire(0, int64(idx)).Right()
		serialization.Ptr2IFace(&h, ptr)
		hosts = append(hosts, h)
	}
	for offset := int32(0); offset < this.backupFactor-1; offset++ {
		ridx := (mod + offset) % (this.remoteNum)
		hosts = append(hosts, this.remoteHosts[ridx])
	}

	var master host.Host
	slaves := make([]host.Host, 0, 0)
	for _, h := range hosts {
		if h.IsValid() {
			if master == nil {
				master = h
			} else {
				slaves = append(slaves, h)
			}
		}
	}

	if master == nil {
		ret.Error(fmt.Errorf("no available master host found for id: %d", id))
		return
	}

	ret.Value(host.NewDuplicatedHost(master, slaves).Right())
	return
}