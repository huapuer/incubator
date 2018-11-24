package network

import (
	"../common/maybe"
	"bufio"
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"time"
)

type Server interface {
	Start(Server, context.Context, string, string) maybe.MaybeError
	handleConnection(Server, context.Context, net.Conn)
	handleData([]byte, int) maybe.MaybeError
	handlePackage([]byte) maybe.MaybeError
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
		server.handleData(buffer, len).Test()
	}
	return
}

//go:noescape
func (this commonServer) handlePacakge(server Server, data []byte) (err maybe.MaybeError) {
	server.handlePackage(data)
	return
}
