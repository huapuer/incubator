package io

import (
	"fmt"
	"github.com/incubator/common/maybe"
	"github.com/incubator/config"
	"github.com/incubator/message"
	"github.com/incubator/network"
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

type joint struct {
	begin  int64
	end    int64
	client network.Client
}

type IO interface {
	config.IOC

	SetLayer(int32)
	Input(int64, message.RemoteMessage) maybe.MaybeError
	Output(int64, message.RemoteMessage) maybe.MaybeError
}

type MaybeIO struct {
	config.IOC

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

type commonIO struct {
	layerId int32
}

func (this *commonIO) SetLayer(id int32) {
	this.layerId = id
}
