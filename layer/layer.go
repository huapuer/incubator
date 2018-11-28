package layer

import (
	"../common/maybe"
	"../config"
	"../message"
	"../router"
	"errors"
	"fmt"
	"../host"
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
	return
}

func SetLayer(cfg config.Config) (err maybe.MaybeError) {
	if _, ok := layers[cfg.Layer.Id]; ok {
		err.Error(fmt.Errorf("layer has been set: %d", cfg.Layer.Id))
		return
	}
	if prototype, ok := layerPrototype[cfg.Layer.Class]; ok {
		layer := prototype.New(cfg.Layer.Attributes, cfg).(MaybeLayer).Right()
		layers[cfg.Layer.Id] = layer
		return
	}
	err.Error(fmt.Errorf("layer prototype not found: %s", cfg.Layer.Class))
	return
}

func GetLayer(layer int32) (ret MaybeLayer) {
	if layer, ok := layers[layer]; ok {
		ret.Value(layer)
		return
	}
	ret.Error(fmt.Errorf("global layer not found: %d", layer))
	return
}

type Layer interface {
	config.IOC

	GetRouter(int32) router.MaybeRouter
	GetMessageFromClass(string) message.MaybeRemoteMessage
	GetMessageCanonicalFromType(int32) message.MaybeRemoteMessage
	LookupHost(int64) host.MaybeHost
	LookupLink(int64, int64) host.MaybeHost
}

type MaybeLayer struct {
	config.IOC

	maybe.MaybeError
	value Layer
}

func (this MaybeLayer) New(cfg config.Config) config.IOC {
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

type commonLayer struct {
	space                     string
	layer                     int32
	messageCanonicalFromClass map[string]message.RemoteMessage
	messageCanonicalFromType  map[int32]message.RemoteMessage
	routers                   map[int32]router.Router
	messageRouters            map[int32]router.Router
}

func (this *commonLayer) Init(attrs interface{}, cfg config.Config) (err maybe.MaybeError) {
	if cfg.Layer.Space == "" {
		err.Error(errors.New("layer space not set"))
		return
	}
	if cfg.Layer.Id < 0 {
		err.Error(fmt.Errorf("illegal layer id: %d", cfg.Layer.Id))
		return
	}

	if cfg.Layer.Id <= 0 {
		err.Error(fmt.Errorf("illegal layer layer: %d", cfg.Layer.Id))
		return
	}
	if cfg.Layer.Space == "" {
		err.Error(errors.New("empty layer space"))
		return
	}
	this.layer = cfg.Layer.Id
	this.space = cfg.Layer.Space

	this.messageRouters = make(map[int32]router.Router)

	for _, routerCfg := range cfg.Routers {
		if _, ok := this.routers[routerCfg.Id]; ok {
			err.Error(fmt.Errorf("router already exists: %d", routerCfg.Id))
			return
		}
		routerPrototype := router.GetRouterPrototype(routerCfg.Class).Right()
		this.routers[routerCfg.Id] = routerPrototype
	}

	this.messageCanonicalFromClass = make(map[string]message.RemoteMessage)
	this.messageCanonicalFromType = make(map[int32]message.RemoteMessage)

	for _, msgCfg := range cfg.Messages {
		if _, ok := this.messageCanonicalFromType[msgCfg.Type]; ok {
			err.Error(fmt.Errorf("message canonical type already exists: %d", msgCfg.Type))
			return
		}
		if _, ok := this.messageCanonicalFromClass[msgCfg.Class]; ok {
			err.Error(fmt.Errorf("message canonical class already exists: %s", msgCfg.Class))
			return
		}
		if _, ok := this.messageRouters[msgCfg.RouterId]; !ok {
			err.Error(fmt.Errorf("router %d not found when register message type %d", msgCfg.RouterId, msgCfg.Type))
		}

		msgCanon := message.GetMessagePrototype(msgCfg.Class).Right()
		msgCanon.SetLayer(int8(this.layer))
		msgCanon.SetType(int8(msgCfg.Type))

		// TODO: deep copy
		this.messageCanonicalFromType[msgCfg.Type] = msgCanon
		this.messageCanonicalFromClass[msgCfg.Class] = msgCanon

		this.messageRouters[msgCfg.Type], _ = this.routers[msgCfg.RouterId]
	}

	err.Error(nil)
	return
}

func (this commonLayer) GetRouter(id int32) (ret router.MaybeRouter) {
	if val, ok := this.routers[id]; ok {
		ret.Value(val)
		return
	}
	ret.Error(fmt.Errorf("router not found: %d", id))
	return
}

func (this commonLayer) GetMessageFromClass(name string) (ret message.MaybeRemoteMessage) {
	if val, ok := this.messageCanonicalFromClass[name]; ok {
		ret.Value(val)
		return
	}
	ret.Error(fmt.Errorf("message canonical from class not found: %s", name))
	return
}

func (this commonLayer) GetMessageCanonicalFromType(typ int32) (ret message.MaybeRemoteMessage) {
	if val, ok := this.messageCanonicalFromType[typ]; ok {
		ret.Value(val)
		return
	}
	ret.Error(fmt.Errorf("message canonical from type not found: %d", typ))
	return
}
