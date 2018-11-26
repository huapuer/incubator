package link

import (
	"fmt"
	"../common/maybe"
	"../config"
	"../storage"
)

var (
	linkPrototypes = make(map[string]Link)
)

func RegisterLinkPrototype(name string, val Link) (err maybe.MaybeError) {
	if _, ok := linkPrototypes[name]; ok {
		err.Error(fmt.Errorf("link prototype redefined: %s", name))
		return
	}
	linkPrototypes[name] = val
	return
}

func GetLinkPrototype(name string) (ret MaybeLink) {
	if prototype, ok := linkPrototypes[name]; ok {
		ret.Value(prototype)
		return
	}
	ret.Error(fmt.Errorf("link prototype for class not found: %s", name))
	return
}

type Link interface {
	storage.DenseTableElement

	GetToId() int64
	SetToId(int64)
}

type MaybeLink struct {
	config.IOC

	maybe.MaybeError
	value Link
}

func (this MaybeLink) Value(value Link) {
	this.Error(nil)
	this.value = value
}

func (this MaybeLink) Right() Link {
	this.Test()
	return this.value
}

func (this MaybeLink) New(cfg config.Config, args ...int32) config.IOC {
	panic("not implemented.")
}

type commonLink struct {
	toId int64
}

func (this commonLink) GetToId() int64 {
	return this.toId
}

func (this commonLink) SetToId(id int64) {
	this.toId = id
}
