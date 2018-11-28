package network

import (
	"../common/maybe"
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/incubator/protocal"
	"math/rand"
	"net"
	"time"
)

type Server interface {
	Start(Server, context.Context, string, string) maybe.MaybeError
	handleConnection(Server, context.Context, net.Conn)
	handleData(Server, []byte, int, net.Conn) maybe.MaybeError
	handlePackage([]byte, net.Conn) maybe.MaybeError
}

type commonServer struct {
	readBufferSize int
	packageBuffer  []byte
	packageSize    int
	headerSize     int
	p              protocal.Protocal
}

func (this commonServer) Start(server Server, ctx context.Context, network string, port string) (err maybe.MaybeError) {
	this.packageBuffer = make([]byte, this.readBufferSize, 0)

	if network == "" {
		err.Error(errors.New("not network provided"))
	}
	if port == "" {
		err.Error(errors.New("no port provided"))
	}
	port = ":" + port
	l, e := net.Listen(network, port)
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
			go server.handleConnection(server, ctx, c)
		}, nil)
	}
}

func (this commonServer) handleConnection(server Server, ctx context.Context, c net.Conn) {
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
		this.handleData(server, buffer, len, c).Test()
	}
	return
}

//go:noescape
func (this commonServer) handlePacakge(server Server, data []byte, c net.Conn) (err maybe.MaybeError) {
	server.handlePackage(data, c)
	return
}

//go:noescape
func (this commonServer) handleData(server Server, data []byte, l int, c net.Conn) (err maybe.MaybeError) {
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
			server.handlePackage(pkg).Test()
			this.packageSize = protocal.PROTOCAL_PARSE_STATE_SHORT
		}
	} else if this.packageSize == protocal.PROTOCAL_PARSE_STATE_ERROR {
		this.packageSize = protocal.PROTOCAL_PARSE_STATE_SHORT
		this.packageBuffer = make([]byte, this.readBufferSize, 0)
	}

	return
}
