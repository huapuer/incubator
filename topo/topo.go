package topo

import (
	"errors"
	"fmt"
	"../common/maybe"
	"../config"
	"../host"
	"../message"
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

func CheckTopo(layer int32)(err maybe.MaybeError) {
	if _, ok :=topos[layer]; ok{
		err.Error(errors.New("global topo has been set"))
		return
	}
	return
}

func SetTopo(layerOffset int32, layer int32, cfg config.Config) (err maybe.MaybeError) {
	if _, ok :=topos[layer]; ok{
		err.Error(fmt.Errorf("topo has been set: %d", layer))
		return
	}
	if prototype, ok := topoPrototype[cfg.Topo.Class]; ok {
		topo := prototype.New(cfg.Topo.Attributes, cfg).(MaybeTopo).Right()

		for typ, msgCfg := range cfg.Messages {
			msg := message.GetMessageCanonical(layerOffset + typ).Right()
			topo.RegisterMessageCanonical(msgCfg.Class, msg).Test()
		}

		topos[layer] = topo
		return
	}
	err.Error(fmt.Errorf("topo prototype not found: %s", cfg.Topo.Class))
	return
}

func GetTopo(layer int32) (ret MaybeTopo) {
	if topo, ok := topos[layer];ok {
		ret.Value(topo)
		return
	}
	ret.Error(fmt.Errorf("global topo not found: %d", layer))
	return
}

type Topo interface {
	config.IOC

	Lookup(int64) host.MaybeHost
	RegisterMessageCanonical(string, message.Message) maybe.MaybeError
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

func (this *commonTopo) RegisterMessageCanonical(className string, msg message.Message) (err maybe.MaybeError){
	if className == "" {
		err.Error(error.Error("empty class name"))
		return
	}
	if msg == nil {
		err.Error(error.Error("message is nil"))
		return
	}
	if _, ok := this.messageCanonical[className]; ok{
		err.Error(fmt.Errorf("message canonical already exists: %s", className))
		return
	}
	this.messageCanonical[className] = msg
	return
}
