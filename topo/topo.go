package topo

import (
	"../common/maybe"
	"../config"
	"../host"
	"fmt"
	"unsafe"
)

var (
	topoPrototypes = make(map[string]Topo)
)

func RegisterTopoPrototype(name string, val Topo) (err maybe.MaybeError) {
	if _, ok := topoPrototypes[name]; ok {
		err.Error(fmt.Errorf("host prototype redefined: %s", name))
		return
	}
	topoPrototypes[name] = val
	return
}

func GetTopoPrototype(name string) (ret MaybeTopo) {
	if prototype, ok := topoPrototypes[name]; ok {
		ret.Value(prototype)
		return
	}
	ret.Error(fmt.Errorf("host prototype for class not found: %s", name))
	return
}

type Topo interface {
	config.IOC

	LookupHost(int64) host.MaybeHost
	LookupLink(int64, int64) host.MaybeHost
	TraverseLinksOfHost(int64, func(ptr unsafe.Pointer) bool) maybe.MaybeError
	GetRemoteHosts() []host.Host
}

type MaybeTopo struct {
	config.IOC

	maybe.MaybeError
	value Topo
}

func (this MaybeTopo) Value(value Topo) {
	this.Error(nil)
	this.value = value
}

func (this MaybeTopo) Right() Topo {
	this.Test()
	return this.value
}
