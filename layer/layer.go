package layer

import (
	"../common/maybe"
	"../config"
	"../message"
	"../network"
	"../router"
	"../serialization"
	"../topo"
	"errors"
	"fmt"
	"github.com/incubator/io"
	"math/rand"
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

func AddLayer(cfg config.Config) (err maybe.MaybeError) {
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
	config.IOC

	GetRouter(int32) router.MaybeRouter
	GetMessageType(interface{}) maybe.MaybeInt32
	GetMessageCanonicalFromType(int32) message.MaybeRemoteMessage
	Start()
	GetTopo() topo.Topo
	GetServer() network.Server
	Stop()
	GetConfig() *config.Config
	GetVersion() int64
	GetSuperLayer() int32
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

type CommonLayer struct {
	space                    string
	id                       int32
	messageClassToType       map[int]int32
	messageCanonicalFromType map[int32]message.RemoteMessage
	routers                  map[int32]router.Router
	messageRouters           map[int32]router.Router
	server                   network.Server
	cfg                      *config.Config
	version                  int64
	superLayer               int32
	io                       io.IO
}

func (this *CommonLayer) Init(attrs interface{}, cfg config.Config) (err maybe.MaybeError) {
	if cfg.Layer.Space == "" {
		err.Error(errors.New("layer space not set"))
		return
	}
	if cfg.Layer.Id < 0 {
		err.Error(fmt.Errorf("illegal layer id: %d", cfg.Layer.Id))
		return
	}
	if cfg.Layer.SuperLayer < 0 {
		err.Error(fmt.Errorf("illegal supervisor layer id: %d", cfg.Layer.SuperLayer))
		return
	}

	this.id = cfg.Layer.Id
	this.space = cfg.Layer.Space

	this.messageRouters = make(map[int32]router.Router)

	this.messageClassToType = make(map[int]int32)
	this.messageCanonicalFromType = make(map[int32]message.RemoteMessage)

	for _, msgCfg := range cfg.Messages {
		if _, ok := this.messageCanonicalFromType[msgCfg.Type]; ok {
			err.Error(fmt.Errorf("message canonical type already exists: %d", msgCfg.Type))
			return
		}

		msgCanon := message.GetMessagePrototype(msgCfg.Class).Right()
		msgCanon.SetLayer(int8(this.id))
		msgCanon.SetType(int8(msgCfg.Type))

		this.messageCanonicalFromType[msgCfg.Type] = msgCanon
		_type := serialization.Eface2TypeInt(msgCanon)
		this.messageClassToType[_type] = msgCfg.Type

		r, ok := this.routers[msgCfg.RouterId]
		if !ok {
			routerCfg, ok := cfg.Routers[msgCfg.RouterId]
			if !ok {
				err.Error(fmt.Errorf("router cfg for %d not found when register message type %d",
					msgCfg.RouterId, msgCfg.Type))
				return
			}
			this.routers[msgCfg.RouterId] = router.GetRouterPrototype(routerCfg.Class).
				Right().New(routerCfg.Attributes, cfg).(router.MaybeRouter).Right()
		}
		this.messageRouters[msgCfg.Type] = r
	}

	if cfg.IO.Class != "" {
		this.io = io.GetIOPrototype(cfg.IO.Class).Right().New(cfg.IO.Attributes, cfg).(io.MaybeIO).Right()
	}

	this.superLayer = cfg.Layer.SuperLayer
	this.cfg = &cfg
	this.version = rand.Int63()

	err.Error(nil)
	return
}

func (this CommonLayer) GetRouter(id int32) (ret router.MaybeRouter) {
	if val, ok := this.routers[id]; ok {
		ret.Value(val)
		return
	}
	ret.Error(fmt.Errorf("router not found: %d", id))
	return
}

//go:noescape
func (this CommonLayer) GetMessageType(msg interface{}) (ret maybe.MaybeInt32) {
	_type := serialization.Eface2TypeInt(msg)
	if typ, ok := this.messageClassToType[_type]; ok {
		ret.Value(typ)
		return
	}
	ret.Error(fmt.Errorf("message type from not found: %+v", msg))
	return
}

func (this CommonLayer) GetMessageCanonicalFromType(typ int32) (ret message.MaybeRemoteMessage) {
	if val, ok := this.messageCanonicalFromType[typ]; ok {
		ret.Value(val)
		return
	}
	ret.Error(fmt.Errorf("message canonical from type not found: %d", typ))
	return
}

func (this CommonLayer) GetConfig() *config.Config {
	return this.cfg
}

func (this CommonLayer) GetVersion() int64 {
	return this.version
}

func (this CommonLayer) GetSuperLayer() int32 {
	return this.superLayer
}
