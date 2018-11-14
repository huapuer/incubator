package topo

import (
	"errors"
	"fmt"
	"../common/maybe"
	"../config"
	"../host"
	"../message"
	"github.com/incubator/router"
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

func SetTopo(cfg config.Config) (err maybe.MaybeError) {
	if _, ok :=topos[cfg.Topo.Layer]; ok{
		err.Error(fmt.Errorf("topo has been set: %d", cfg.Topo.Layer))
		return
	}
	if prototype, ok := topoPrototype[cfg.Topo.Class]; ok {
		topo := prototype.New(cfg.Topo.Attributes, cfg).(MaybeTopo).Right()
		topos[cfg.Topo.Layer] = topo
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
	GetRouter(int32) router.MaybeRouter
	GetMessageCanonicalFromClass(string) message.MaybeMessage
	GetMessageCanonicalFromType(int32) message.MaybeMessage
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
	layer int32
	messageCanonicalFromClass map[string]message.Message
	messageCanonicalFromType map[int32]message.Message
	routers map[int32]router.Router
	messageRouters map[int32]router.Router
}

func (this *commonTopo) Init(cfg config.Config) (err maybe.MaybeError){
	this.layer = cfg.Topo.Layer

	this.messageRouters = make(map[int32]router.Router)

	for _, routerCfg := range cfg.Routers {
		if _, ok := this.routers[routerCfg.Id]; ok {
			err.Error(fmt.Errorf("router already exists: %d", routerCfg.Id))
			return
		}
		routerPrototype := router.GetRouterPrototype(routerCfg.Class).Right()
		this.routers[routerCfg.Id] = routerPrototype
	}

	this.messageCanonicalFromClass = make(map[string]message.Message)
	this.messageCanonicalFromType = make(map[int32]message.Message)

	for _, msgCfg := range cfg.Messages {
		if _, ok := this.messageCanonicalFromType[msgCfg.Type]; ok {
			err.Error(fmt.Errorf("message canonical type already exists: %d", msgCfg.Type))
			return
		}
		if _, ok := this.messageCanonicalFromClass[msgCfg.Class]; ok{
			err.Error(fmt.Errorf("message canonical class already exists: %s", msgCfg.Class))
			return
		}
		if _, ok := this.messageRouters[msgCfg.RouterId]; !ok {
			err.Error(fmt.Errorf("router %d not found when register message type %d", msgCfg.RouterId, msgCfg.Type))
		}

		msgPrototype := message.GetMessagePrototype(msgCfg.Class).Right()
		msgCanon := msgPrototype.Duplicate().Right()
		msgCanon.SetLayer(this.layer)
		msgCanon.SetType(msgCfg.Type)

		// TODO: deep copy
		this.messageCanonicalFromType[msgCfg.Type] = msgCanon
		this.messageCanonicalFromClass[msgCfg.Class] = msgCanon

		this.messageRouters[msgCfg.Type], _ = this.routers[msgCfg.RouterId]
	}
	return
}

func (this commonTopo) GetRouter(id int32) (ret router.MaybeRouter) {
	if val, ok := this.routers[id]; ok {
		ret.Value(val)
		return
	}
	ret.Error(fmt.Errorf("router not found: %d", id))
	return
}

func (this commonTopo) GetMessageFromClass(name string) (ret message.MaybeMessage) {
	if val, ok := this.messageCanonicalFromClass[name]; ok {
		ret.Value(val.Duplicate().Right())
		return
	}
	ret.Error(fmt.Errorf("message canonical from class not found: %s", name))
	return
}

func (this commonTopo) GetMessageCanonicalFromType(typ int32) (ret message.MaybeMessage) {
	if val, ok := this.messageCanonicalFromType[typ]; ok {
		ret.Value(val)
		return
	}
	ret.Error(fmt.Errorf("message canonical from type not found: %d", typ))
	return
}
