package interfaces

import (
	"fmt"
	"github.com/incubator/common/maybe"
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
	IOC

	SendToHost(int64, RemoteMessage) maybe.MaybeError
	SendToLink(int64, int64, RemoteMessage) maybe.MaybeError
	TraverseOutLinksOfHost(int64, func(ptr unsafe.Pointer) bool) maybe.MaybeError
	GetRemoteHosts() []Host
	LookupHost(int64) MaybeHost
	Start()
	GetLayer() int32
	SetLayer(int32)
	GetRemoteHostId(int32) int64
	GetAddr() string
}

type MaybeTopo struct {
	IOC

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

func (this MaybeTopo) New(attr interface{}, cfg Config) IOC {
	panic("not implemented.")
}
