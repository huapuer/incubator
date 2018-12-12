package interfaces

import (
	"fmt"
	"github.com/incubator/common/maybe"
)

var (
	layerPrototype = make(map[string]Layer)
	layers         = make(map[int32]Layer)
)

func RegisterLayerPrototype(name string, val Layer) (err maybe.MaybeError) {
	if _, ok := layerPrototype[name]; ok {
		err.Error(fmt.Errorf("layer prototype redefined: %s", name))
		return
	}
	layerPrototype[name] = val

	err.Error(nil)
	return
}

func AddLayer(cfg Config) (err maybe.MaybeError) {
	if _, ok := layers[cfg.GetLayerId()]; ok {
		err.Error(fmt.Errorf("layer has been set: %d", cfg.GetLayerId()))
		return
	}
	if prototype, ok := layerPrototype[cfg.GetLayerClass()]; ok {
		layer := prototype.New(cfg.GetLayerAttr(), cfg).(MaybeLayer).Right()
		layers[cfg.GetLayerId()] = layer
		return
	}
	err.Error(fmt.Errorf("layer prototype not found: %s", cfg.GetLayerClass()))
	return
}

func DeleteLayer(id int32) (err maybe.MaybeError) {
	layer, ok := layers[id]
	if !ok {
		err.Error(fmt.Errorf("global layer not found: %d", id))
		return
	}
	layer.Stop()
	delete(layers, id)

	err.Error(nil)
	return
}

func GetLayer(id int32) (ret MaybeLayer) {
	if layer, ok := layers[id]; ok {
		ret.Value(layer)
		return
	}
	ret.Error(fmt.Errorf("global layer not found: %d", id))
	return
}

type Layer interface {
	IOC

	GetRouter(int32) MaybeRouter
	GetMessageType(interface{}) maybe.MaybeInt32
	GetMessageCanonicalFromType(int32) MaybeRemoteMessage
	Start()
	GetTopo() Topo
	GetService(int32) MaybeServer
	Stop()
	GetConfig() Config
	GetVersion() int64
	GetSuperLayer() int32
}

type MaybeLayer struct {
	IOC

	maybe.MaybeError
	value Layer
}

func (this MaybeLayer) New(attr interface{}, cfg Config) IOC {
	panic("not implemented.")
}

func (this MaybeLayer) Value(value Layer) {
	this.Error(nil)
	this.value = value
}

func (this MaybeLayer) Right() Layer {
	this.Test()
	return this.value
}
