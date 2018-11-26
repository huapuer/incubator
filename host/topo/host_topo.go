package topo

import (
	"../../host"
	"incubator/config"
	"incubator/common/maybe"
	"fmt"
)

var (
	hostTopoPrototypes = make(map[string]HostTopo)
)

func RegisterHostTopoPrototype(name string, val HostTopo) (err maybe.MaybeError) {
	if _, ok := hostTopoPrototypes[name]; ok {
		err.Error(fmt.Errorf("host prototype redefined: %s", name))
		return
	}
	hostTopoPrototypes[name] = val
	return
}

func GetHostTopoPrototype(name string) (ret MaybeHostTopo) {
	if prototype, ok := hostTopoPrototypes[name]; ok {
		ret.Value(prototype)
		return
	}
	ret.Error(fmt.Errorf("host prototype for class not found: %s", name))
	return
}

type HostTopo interface {
	config.IOC

	Lookup(id int64) (ret host.MaybeHost)
}

type MaybeHostTopo struct {
	config.IOC

	maybe.MaybeError
	value HostTopo
}

func (this MaybeHostTopo) Value(value HostTopo) {
	this.Error(nil)
	this.value = value
}

func (this MaybeHostTopo) Right() HostTopo {
	this.Test()
	return this.value
}