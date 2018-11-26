package layer

import (
	"../config"
	"errors"
	"fmt"
	"github.com/incubator/host"
	"github.com/incubator/topo"
)

const (
	defaultLayerClassName = "layer.defaultLayer"
)

func init() {
	RegisterLayerPrototype(defaultLayerClassName, &defaultLayer{})
}

type defaultLayer struct {
	commonLayer

	topo topo.Topo
}

func (this *defaultLayer) New(attrs interface{}, cfg config.Config) config.IOC {
	ret := MaybeLayer{}

	layer := &defaultLayer{
		commonLayer: commonLayer{},
	}

	layer.Init(attrs, cfg).Test()

	attrsMap, ok := cfg.Layer.Attributes.(map[string]interface{})
	if !ok {
		ret.Error(fmt.Errorf("layer attrs cfg type error(expecting map[stirng]interface{}): %+v", cfg.Layer.Attributes))
	}

	topoSchema, ok := attrsMap["TopoSchema"]
	if !ok {
		ret.Error(errors.New("topo schema not set"))
		return ret
	}

	topoSchemaInt, ok := topoSchema.(int32)
	if !ok {
		ret.Error(fmt.Errorf("topo class cfg type error(expecting int32): %+v", topoSchema))
		return ret
	}
	if topoSchemaInt <= 0 {
		ret.Error(fmt.Errorf("illegal TopoSchema: %d", topoSchemaInt))
		return ret
	}

	topoCfg, ok := cfg.Topos[topoSchemaInt]
	if !ok {
		ret.Error(fmt.Errorf("topo cfg not found: %d", topoSchemaInt))
		return ret
	}

	layer.topo = topo.GetTopoPrototype(topoCfg.Class).Right().(config.IOC).New(topoCfg.Attributes, cfg).(topo.Topo)

	ret.Value(layer)
	return ret
}

func (this defaultLayer) LookupHost(id int64) host.MaybeHost {
	return this.topo.LookupHost(id)
}

func (this defaultLayer) LookupLink(hid int64, gid int64) host.MaybeHost {
	return this.topo.LookupLink(hid, gid)
}
