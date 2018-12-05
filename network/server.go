package network

import (
	"../common/class"
	"../common/maybe"
	"../config"
	"../protocal"
	"bufio"
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"time"
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
	config.IOC

	Start(context.Context) maybe.MaybeError
	HandleConnection(context.Context, net.Conn)
	handleData([]byte, int, net.Conn) maybe.MaybeError
	handlePackage([]byte, net.Conn) maybe.MaybeError
}

type MaybeServer struct {
	config.IOC

	maybe.MaybeError
	value Server
}

func (this MaybeServer) Value(value Server) {
	this.Error(nil)
	this.value = value
}

func (this MaybeServer) Right() Server {
	this.Test()
	return this.value
}

func (this MaybeServer) New(cfg config.Config, args ...int32) config.IOC {
	panic("not implemented.")
}

type commonServer struct {
	class.Class

	network        string
	address        string
	readBufferSize int
	packageBuffer  []byte
	packageSize    int
	headerSize     int
	p              protocal.Protocal
	derived        Server
}

func (this commonServer) Start(ctx context.Context) (err maybe.MaybeError) {
	this.packageBuffer = make([]byte, this.readBufferSize, 0)

	if this.network == "" {
		err.Error(errors.New("not network provided"))
	}
	if this.address == "" {
		err.Error(errors.New("no address provided"))
	}
	l, e := net.Listen(this.network, this.address)
	if e != nil {
		err.Error(e)
		return
	}
	defer l.Close()
	rand.Seed(time.Now().Unix())

	for {
		select {
		case <-ctx.Done():
			return
		}
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		maybe.TryCatch(func() {
			go this.HandleConnection(ctx, c)
		}, nil)
	}
}

func (this commonServer) HandleConnection(ctx context.Context, c net.Conn) {
	defer c.Close()
	reader := bufio.NewReader(c)
	buffer := make([]byte, this.readBufferSize)
	for {
		select {
		case <-ctx.Done():
			return
		}
		len, err := reader.Read(buffer)
		if err != nil {
			panic(err)
		}
		this.handleData(buffer, len, c).Test()
	}
	return
}

//go:noescape
func (this commonServer) handlePacakge(data []byte, c net.Conn) (err maybe.MaybeError) {
	this.derived.handlePackage(data, c)
	return
}

//go:noescape
func (this commonServer) handleData(data []byte, l int, c net.Conn) (err maybe.MaybeError) {
	if l == 0 {
		err.Error(errors.New("empty data"))
		return
	}
	this.packageBuffer = append(this.packageBuffer, data...)
	if this.packageSize == protocal.PROTOCAL_PARSE_STATE_SHORT {
		this.packageSize, this.headerSize = this.p.Parse(data)
	}
	if this.packageSize >= 0 {
		if this.readBufferSize >= this.packageSize {
			pkg := this.packageBuffer[this.headerSize:this.packageSize]
			this.packageBuffer = this.packageBuffer[this.packageSize:]
			this.derived.handlePackage(pkg, c).Test()
			this.packageSize = protocal.PROTOCAL_PARSE_STATE_SHORT
		}
	} else if this.packageSize == protocal.PROTOCAL_PARSE_STATE_ERROR {
		this.packageSize = protocal.PROTOCAL_PARSE_STATE_SHORT
		this.packageBuffer = make([]byte, this.readBufferSize, 0)
	}

	return
}

func (this *commonServer) Inherit(cobj class.Class) {
	this.derived = cobj.(Server)
}
