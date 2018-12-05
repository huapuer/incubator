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

	if cfg.Server.Class != "" {
		this.server = network.GetServerPrototype(cfg.Server.Class).Right().New(nil, cfg).(network.Server)
	}

	ret.Value(layer)
	return ret
}

func (this defaultLayer) Start() {
	for _, r := range this.routers {
		r.Start()
	}
	this.server.Start(context.Background()).Test()
	this.topo.SetLayer(this.id)
	this.topo.Start()
}

func (this defaultLayer) GetTopo() topo.Topo {
	return this.topo
}

func (this defaultLayer) GetServer() network.Server {
	return this.server
}
