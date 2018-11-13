package router

import (
	"fmt"
	"../common/maybe"
	"../config"
	"../message"
)

var (
	routerPrototype = make(map[string]Router)
	routers         = make(map[int32]Router)
)

func RegisterRouterPrototype(name string, val Router) (err maybe.MaybeError) {
	if _, ok := routerPrototype[name]; ok {
		err.Error(fmt.Errorf("router redefined: %s", name))
		return
	}
	routerPrototype[name] = val
	return
}

func AddRouter(layerOffset int32, id int32, className string, cfg config.Config) (err maybe.MaybeError) {
	if _, ok := routers[id]; ok {
		err.Error(fmt.Errorf("router already exists: %d", id))
		return
	}
	if prototype, ok := routerPrototype[className]; ok {
		routerCfg, ok := cfg.Routers[id]
		if !ok {
			err.Error(fmt.Errorf("router cfg not found: %s", className))
			return
		}
		newRouter := prototype.New(routerCfg.Attributes, cfg).(MaybeRouter).Right()
		routers[layerOffset + id] = newRouter
		return
	}
	err.Error(fmt.Errorf("router prototype not found: %s", className))
	return
}

func GetRouter(id int32) (ret MaybeRouter) {
	if val, ok := routers[id]; ok {
		ret.Value(val)
		return
	}
	ret.Error(fmt.Errorf("router not found: %d", id))
	return
}

type Router interface {
	config.IOC

	Route(message.Message) maybe.MaybeError
}

type MaybeRouter struct {
	maybe.MaybeError
	value Router
}

func (this MaybeRouter) New(attrs interface{}, cfg config.Config) config.IOC {
	panic("not implemented.")
}

func (this MaybeRouter) Value(value Router) {
	this.Error(nil)
	this.value = value
}

func (this MaybeRouter) Right() Router {
	this.Test()
	return this.value
}
