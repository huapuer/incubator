package interfaces

import (
	"context"
	"fmt"
	"github.com/incubator/common/maybe"
	"net"
)

var (
	serverPrototypes = make(map[string]Server)
)

func RegisterServerPrototype(name string, val Server) (err maybe.MaybeError) {
	if _, ok := serverPrototypes[name]; ok {
		err.Error(fmt.Errorf("server prototype redefined: %s", name))
		return
	}
	serverPrototypes[name] = val

	err.Error(nil)
	return
}

func GetServerPrototype(name string) (ret MaybeServer) {
	if prototype, ok := serverPrototypes[name]; ok {
		ret.Value(prototype)
		return
	}
	ret.Error(fmt.Errorf("server prototype for class not found: %s", name))
	return
}

type Server interface {
	IOC

	Start(context.Context) maybe.MaybeError
	HandleConnection(context.Context, net.Conn)
	HandleData([]byte, int, net.Conn) maybe.MaybeError
	HandlePackage([]byte, net.Conn) maybe.MaybeError
	SetPort(int)
	GetProtocal() Protocal
}

type MaybeServer struct {
	IOC

	maybe.MaybeError
	value Server
}

func (this *MaybeServer) Value(value Server) {
	this.Error(nil)
	this.value = value
}

func (this MaybeServer) Right() Server {
	this.Test()
	return this.value
}

func (this MaybeServer) New(attr interface{}, cfg Config) IOC {
	panic("not implemented.")
}

type Client interface {
	IOC

	Connect(string)
	Send(message RemoteMessage) maybe.MaybeError
}

type MaybeClient struct {
	IOC

	maybe.MaybeError
	value Client
}

func (this *MaybeClient) Value(value Client) {
	this.Error(nil)
	this.value = value
}

func (this MaybeClient) Right() Client {
	this.Test()
	return this.value
}

func (this MaybeClient) New(attrs interface{}, cfg Config) IOC {
	panic("not implemented.")
}
