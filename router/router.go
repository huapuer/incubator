package router

import (
	"../common/maybe"
	"../config"
	"../message"
	"fmt"
)

var (
	routerPrototype = make(map[string]Router)
)

func RegisterRouterPrototype(name string, val Router) (err maybe.MaybeError) {
	if _, ok := routerPrototype[name]; ok {
		err.Error(fmt.Errorf("router redefined: %s", name))
		return
	}
	routerPrototype[name] = val
	return
}

func GetRouterPrototype(name string) (ret MaybeRouter) {
	if routerPrototype, ok := routerPrototype[name]; ok {
		ret.Value(routerPrototype)
		return
	}
	ret.Error(fmt.Errorf("router prototype not found: %s", name))
	return
}

type Router interface {
	config.IOC

	Start()
	Route(message.RemoteMessage) maybe.MaybeError
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
