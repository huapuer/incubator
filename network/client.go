package network

import (
	"../common/maybe"
	"../config"
	"../message"
	"fmt"
)

var (
	clientPrototypes = make(map[string]Client)
)

func RegisterClientPrototype(name string, val Client) (err maybe.MaybeError) {
	if _, ok := clientPrototypes[name]; ok {
		err.Error(fmt.Errorf("client prototype redefined: %s", name))
		return
	}
	clientPrototypes[name] = val
	return
}

func GetClientPrototype(name string) (ret MaybeClient) {
	if prototype, ok := clientPrototypes[name]; ok {
		ret.Value(prototype)
		return
	}
	ret.Error(fmt.Errorf("client prototype for class not found: %s", name))
	return
}

type Client interface {
	config.IOC

	Connect(string)
	Send(message message.RemoteMessage) maybe.MaybeError
}

type MaybeClient struct {
	config.IOC

	maybe.MaybeError
	value Client
}

func (this MaybeClient) Value(value Client) {
	this.Error(nil)
	this.value = value
}

func (this MaybeClient) Right() Client {
	this.Test()
	return this.value
}

func (this MaybeClient) New(attrs interface{}, cfg config.Config) config.IOC {
	panic("not implemented.")
}
