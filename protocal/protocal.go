package protocal

import (
	"fmt"
	"github.com/incubator/common/maybe"
	"github.com/incubator/interfaces"
)

const (
	PROTOCAL_PARSE_STATE_SHORT = -1
	PROTOCAL_PARSE_STATE_ERROR = -2
)

var (
	protocalPrototypes = make(map[string]interfaces.Protocal)
)

func RegisterProtocalPrototype(name string, val interfaces.Protocal) (err maybe.MaybeError) {
	if _, ok := protocalPrototypes[name]; ok {
		err.Error(fmt.Errorf("protocal prototype redefined: %s", name))
		return
	}
	protocalPrototypes[name] = val

	err.Error(nil)
	return
}

func GetProtocalPrototype(name string) (ret interfaces.MaybeProtocal) {
	if prototype, ok := protocalPrototypes[name]; ok {
		ret.Value(prototype)
		return
	}
	ret.Error(fmt.Errorf("protocal prototype for class not found: %s", name))
	return
}
