package topo

import (
	"errors"
	"fmt"
	"../common/maybe"
	"../config"
	"../host"
	"incubator/message"
)

var (
	topoPrototype = make(map[string]Topo)
	topos    = make(map[int32]Topo)
)

func RegisterTopoPrototype(name string, val Topo) (err maybe.MaybeError) {
	if _, ok := topoPrototype[name]; ok {
		err.Error(fmt.Errorf("topo prototype redefined: %s", name))
		return
	}
	topoPrototype[name] = val
	return
}

func SetTopo(index int32, cfg config.Config) (err maybe.MaybeError) {
	if _, ok :=topos[index]; ok{
		err.Error(errors.New("global topo has been set"))
		return
	}
	if prototype, ok := topoPrototype[cfg.Topo.Class]; ok {
		topo := prototype.New(cfg.Topo.Attributes, cfg).(MaybeTopo).Right()
		topos[index] = topo
		return
	}
	err.Error(fmt.Errorf("topo prototype not found: %s", cfg.Topo.Class))
	return
}

func GetTopo(index int32) (ret MaybeTopo) {
	if topo, ok := topos[index];ok {
		ret.Value(topo)
		return
	}
	ret.Error(fmt.Errorf("global topo not found: %d", index))
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

type commonTopo struct{
	messageCanonical map[string]message.Message
}

func (this * commonTopo) RegisterMessageCanonical(name string, val message.Message) (err maybe.MaybeError){
	if _, ok := messagePrototype[name]; ok {
		err.Error(fmt.Errorf("message prototype redefined: %s", name))
		return
	}
	messagePrototype[name] = val
	return
}
