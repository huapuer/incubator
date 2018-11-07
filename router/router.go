package router

import (
	"fmt"
	"../common/maybe"
	"../config"
	"../message"
)

var (
	routerPrototype = make(map[string]Router)
	routers         = make(map[string]Router)
)

func RegisterRouterPrototype(name string, val Router) (err maybe.MaybeError) {
	if _, ok := routerPrototype[name]; ok {
		err.Error(fmt.Errorf("router redefined: %s", name))
		return
	}
	routerPrototype[name] = val
	return
}

func AddRouter(cfg config.Config, name string) (err maybe.MaybeError) {
	if _, ok := routers[name]; ok {
		err.Error(fmt.Errorf("router already exists: %s", name))
		return
	}
	if prototype, ok := routerPrototype[name]; ok {
		newRouter := prototype.New(cfg).(MaybeRouter).Right()
		routers[name] = newRouter
		return
	}
	err.Error(fmt.Errorf("router prototype not found: %s", name))
	return
}

func GetRouter(name string) (ret MaybeRouter) {
	if val, ok := routers[name]; ok {
		ret.Value(val)
		return
	}
	ret.Error(fmt.Errorf("router not found: %s", name))
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

func (this MaybeRouter) New(cfg config.Config) config.IOC {
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
