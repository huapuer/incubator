package interfaces

import (
	"fmt"
	"github.com/incubator/common/maybe"
)

var (
	ioPrototypes = make(map[string]IO)
)

func RegisterIOPrototype(name string, val IO) (err maybe.MaybeError) {
	if _, ok := ioPrototypes[name]; ok {
		err.Error(fmt.Errorf("io prototype redefined: %s", name))
		return
	}
	ioPrototypes[name] = val
	return
}

func GetIOPrototype(name string) (ret MaybeIO) {
	if prototype, ok := ioPrototypes[name]; ok {
		ret.Value(prototype)
		return
	}
	ret.Error(fmt.Errorf("io prototype for class not found: %s", name))
	return
}

type IO interface {
	IOC

	SetLayer(int32)
	Input(int64, RemoteMessage) maybe.MaybeError
	Output(int64, RemoteMessage) maybe.MaybeError
}

type MaybeIO struct {
	IOC

	maybe.MaybeError
	value IO
}

func (this MaybeIO) Value(value IO) {
	this.Error(nil)
	this.value = value
}

func (this MaybeIO) Right() IO {
	this.Test()
	return this.value
}

func (this MaybeIO) New(attr interface{}, cfg Config) IOC {
	panic("not implemented.")
}
