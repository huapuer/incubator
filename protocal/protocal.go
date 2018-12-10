package protocal

import (
	"../common/maybe"
	"../config"
	"../message"
	"fmt"
)

const (
	PROTOCAL_PARSE_STATE_SHORT = -1
	PROTOCAL_PARSE_STATE_ERROR = -2
)

var (
	protocalPrototypes = make(map[string]Protocal)
)

func RegisterProtocalPrototype(name string, val Protocal) (err maybe.MaybeError) {
	if _, ok := protocalPrototypes[name]; ok {
		err.Error(fmt.Errorf("protocal prototype redefined: %s", name))
		return
	}
	protocalPrototypes[name] = val
	return
}

func GetProtocalPrototype(name string) (ret MaybeProtocal) {
	if prototype, ok := protocalPrototypes[name]; ok {
		ret.Value(prototype)
		return
	}
	ret.Error(fmt.Errorf("protocal prototype for class not found: %s", name))
	return
}

type Protocal interface {
	config.IOC

	Pack(message.RemoteMessage) []byte
	Parse([]byte) (int, int)
	Decode([]byte) []byte
}

type MaybeProtocal struct {
	config.IOC

	maybe.MaybeError
	value Protocal
}

func (this MaybeProtocal) Value(value Protocal) {
	this.Error(nil)
	this.value = value
}

func (this MaybeProtocal) Right() Protocal {
	this.Test()
	return this.value
}

func (this MaybeProtocal) New(cfg config.Config, args ...int32) config.IOC {
	panic("not implemented.")
}
