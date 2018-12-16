package layer

import (
	"context"
	"fmt"
	"github.com/incubator/config"
	"github.com/incubator/interfaces"
)

const (
	defaultLayerClassName = "layer.defaultLayer"
)

func init() {
	interfaces.RegisterLayerPrototype(defaultLayerClassName, &defaultLayer{})
}

type defaultLayer struct {
	CommonLayer

	topo interfaces.Topo
}

func (this *defaultLayer) New(attrs interface{}, cfg interfaces.Config) interfaces.IOC {
	ret := interfaces.MaybeLayer{}

	layer := &defaultLayer{
		CommonLayer: CommonLayer{},
	}

	layer.Init(attrs, cfg.(*config.Config)).Test()

	topoSchema := config.GetAttrInt32(attrs, "TopoSchema", config.CheckInt32GT0).Right()

	topoCfg, ok := cfg.(*config.Config).TopoMap[topoSchema]
	if !ok {
		ret.Error(fmt.Errorf("topo cfg not found: %d", topoSchema))
		return ret
	}

	layer.topo = interfaces.GetTopoPrototype(topoCfg.Class).Right().New(topoCfg.Attributes, cfg).(interfaces.MaybeTopo).Right()

	layer.services = make([]interfaces.Server, 0, 0)
	for _, service := range cfg.(*config.Config).Services {
		serverCfg, ok := cfg.(*config.Config).ServerMap[service.ServerSchema]
		if !ok {
			ret.Error(fmt.Errorf("server config not found for schema: %d", service.ServerSchema))
			return ret
		}
		server := interfaces.GetServerPrototype(serverCfg.Class).Right().New(serverCfg.Attributes, cfg).(interfaces.MaybeServer).Right()
		server.SetPort(service.Port)
		layer.services = append(layer.services, server)
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

func (this defaultLayer) GetTopo() interfaces.Topo {
	return this.topo
}

func (this defaultLayer) GetService(idx int32) (ret interfaces.MaybeServer) {
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
