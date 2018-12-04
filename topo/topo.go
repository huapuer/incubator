package topo

import (
	"../common/maybe"
	"../config"
	"../host"
	"fmt"
	"../message"
	"net"
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

	SendToHost(int64, message.RemoteMessage) maybe.MaybeError
	SendToLink(int64, int64, message.RemoteMessage) maybe.MaybeError
	TraverseOutLinksOfHost(int64, func(ptr unsafe.Pointer) bool) maybe.MaybeError
	GetRemoteHosts() []host.Host
	LookupHost(int64) host.MaybeHost
	GetRemoteHostId(int32) int64
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

type SessionTopo interface {
	Topo

	AddHost(int64, net.Conn) maybe.MaybeError
}
