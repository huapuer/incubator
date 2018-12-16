package network

import (
	"fmt"
	"github.com/incubator/common/maybe"
	"github.com/incubator/interfaces"
)

var (
	clientPrototypes = make(map[string]interfaces.Client)
)

func RegisterClientPrototype(name string, val interfaces.Client) (err maybe.MaybeError) {
	if _, ok := clientPrototypes[name]; ok {
		err.Error(fmt.Errorf("client prototype redefined: %s", name))
		return
	}
	clientPrototypes[name] = val

	err.Error(nil)
	return
}

func GetClientPrototype(name string) (ret interfaces.MaybeClient) {
	if prototype, ok := clientPrototypes[name]; ok {
		ret.Value(prototype)
		return
	}
	ret.Error(fmt.Errorf("client prototype for class not found: %s", name))
	return
}
