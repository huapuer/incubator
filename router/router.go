package router

import (
	"fmt"
	"github.com/incubator/common/maybe"
	"github.com/incubator/interfaces"
)

var (
	routerPrototype = make(map[string]interfaces.Router)
)

func RegisterRouterPrototype(name string, val interfaces.Router) (err maybe.MaybeError) {
	if _, ok := routerPrototype[name]; ok {
		err.Error(fmt.Errorf("router redefined: %s", name))
		return
	}
	routerPrototype[name] = val

	err.Error(nil)
	return
}

func GetRouterPrototype(name string) (ret interfaces.MaybeRouter) {
	if routerPrototype, ok := routerPrototype[name]; ok {
		ret.Value(routerPrototype)
		return
	}
	ret.Error(fmt.Errorf("router prototype not found: %s", name))
	return
}
