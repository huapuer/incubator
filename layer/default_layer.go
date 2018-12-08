package layer

import (
	"../config"
	"../network"
	"../topo"
	"context"
	"fmt"
)

const (
	defaultLayerClassName = "layer.defaultLayer"
)

func init() {
	RegisterLayerPrototype(defaultLayerClassName, &defaultLayer{})
}

type defaultLayer struct {
	CommonLayer

	topo topo.Topo
}

func (this *defaultLayer) New(attrs interface{}, cfg config.Config) config.IOC {
	ret := MaybeLayer{}

	layer := &defaultLayer{
		CommonLayer: CommonLayer{},
	}

	layer.Init(attrs, cfg).Test()

	topoSchema := config.GetAttrInt32(attrs, "TopoSchema", config.CheckInt32GT0).Right()

	topoCfg, ok := cfg.Topos[topoSchema]
	if !ok {
		ret.Error(fmt.Errorf("topo cfg not found: %d", topoSchema))
		return ret
	}

	layer.topo = topo.GetTopoPrototype(topoCfg.Class).Right().New(topoCfg.Attributes, cfg).(topo.Topo)

	this.services = make([]network.Server, 0, 0)
	for _, service := range cfg.Services {
		serverCfg, ok := cfg.Servers[service.ServerSchema]
		if !ok {
			ret.Error(fmt.Errorf("server config not found for schema: %d", service.ServerSchema))
			return ret
		}
		this.services = append(this.services, network.GetServerPrototype(serverCfg.Class).Right().New(serverCfg.Attributes, cfg).(network.Server))
	}

	ret.Value(layer)
	return ret
}

func (this defaultLayer) Start() {
	for _, r := range this.routers {
		r.Start()
	}
	for _, s := range this.services {
		s.Start(context.Background()).Test()
	}
	this.topo.SetLayer(this.id)
	this.topo.Start()
}

func (this defaultLayer) GetTopo() topo.Topo {
	return this.topo
}

func (this defaultLayer) GetService(idx int32) (ret network.MaybeServer) {
	if idx < 0 || int(idx) > len(this.services) {
		ret.Error(fmt.Errorf("service not found: %d", idx))
		return
	}
	ret.Value(this.services[idx])
	return
}

func (this defaultLayer) Stop() {
	for _, r := range this.routers {
		r.Stop()
	}
}
