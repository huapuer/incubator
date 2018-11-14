package network

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"../common/maybe"
	"math/rand"
	"net"
	"time"
)

type Server interface {
	Start(Server, context.Context, string, string) maybe.MaybeError
	handleConnection(Server, context.Context, net.Conn)
	handleData(Server, []byte, int) maybe.MaybeError
	handlePackage(Server, []byte) maybe.MaybeError
}

type commonServer struct {
	readBufferSize int
	packageBuffer  []byte
	packageSize    int
}

func (this commonServer) Start(server Server, ctx context.Context, network string, port string) (err maybe.MaybeError) {
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
		case c, err := l.Accept():
			if err != nil {
				fmt.Println(err)
				return
			}
			maybe.TryCatch(func() {
				go server.handleConnection(server, ctx, c)
			}, nil)
		}
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
		case len, err := reader.Read(buffer):
			if err != nil {
				panic(err)
			}
			server.handleData(server, buffer, len).Test()
		}
	}
	return
}

func (this commonServer) handleData(server Server, data []byte, l int) (err maybe.MaybeError) {
	if l == 0 {
		err.Error(errors.New("empty data"))
		return
	}
	if this.packageSize == 0 {
		this.packageSize = int(data[0])
		this.packageBuffer = data[1:]
	} else {
		want := len(this.packageBuffer) + l - this.packageSize
		if want >= 0 {
			pkg := this.packageBuffer
			pkg = append(pkg, data[:want]...)
			if want > 0 {
				this.packageBuffer = data[want:]
			}
			this.packageSize = 0
			server.handlePackage(server, pkg).Test()
		} else {
			this.packageBuffer = append(this.packageBuffer, data...)
		}
	}

	return
}

func (this commonServer) handlePacakge(server Server, data []byte) (err maybe.MaybeError) {
	server.handlePackage(server, data)
	return
}
