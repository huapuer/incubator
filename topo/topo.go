package topo

import (
	"errors"
	"fmt"
	"../common/maybe"
	"../config"
	"../host"
)

var (
	topoPrototype = make(map[string]Topo)
	globalTopo    Topo
)

func RegisterTopoPrototype(name string, val Topo) (err maybe.MaybeError) {
	if _, ok := topoPrototype[name]; ok {
		err.Error(fmt.Errorf("topo prototype redefined: %s", name))
		return
	}
	topoPrototype[name] = val
	return
}

func SetGlobalTopo(cfg config.Config) (err maybe.MaybeError) {
	if globalTopo != nil {
		err.Error(errors.New("global topo has been set"))
		return
	}
	if prototype, ok := topoPrototype[cfg.Topo.Class]; ok {
		topo := prototype.New(cfg).(MaybeTopo).Right()
		globalTopo = topo
		return
	}
	err.Error(fmt.Errorf("topo prototype not found: %s", cfg.Topo.Class))
	return
}

func GetGlobalTopo() (ret MaybeTopo) {
	if globalTopo == nil {
		ret.Error(errors.New("global topo not set"))
		return
	}
	ret.Value(globalTopo)
	return
}

type Topo interface {
	config.IOC

	Lookup(int64) host.MaybeHost
}

type MaybeTopo struct {
	config.IOC

	maybe.MaybeError
	value Topo
}

func (this MaybeTopo) New(cfg config.Config) config.IOC {
	panic("not implemented.")
}

func (this MaybeTopo) Value(value Topo) {
	this.Error(nil)
	this.value = value
}

func (this MaybeTopo) Right() Topo {
	this.Test()
	return this.value
}
