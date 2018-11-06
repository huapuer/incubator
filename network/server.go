package network

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"../common/class"
	"../common/maybe"
	"math/rand"
	"net"
	"time"
)

type Server interface {
	Start(context.Context, string, string) maybe.MaybeError
	handleConnection(context.Context, net.Conn)
	handleData([]byte, int) maybe.MaybeError
	handlePackage([]byte) maybe.MaybeError
}

type commonServer struct {
	class.DefaultClass

	readBufferSize int
	packageBuffer  []byte
	packageSize    int
}

func (this commonServer) Start(ctx context.Context, network string, port string) (err maybe.MaybeError) {
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
				go this.GetDerived().(Server).handleConnection(ctx, c)
			}, nil)
		}
	}
}

func (this commonServer) handleConnection(ctx context.Context, c net.Conn) {
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
			this.GetDerived().(Server).handleData(buffer, len).Test()
		}
	}
	return
}

func (this commonServer) handleData(data []byte, l int) (err maybe.MaybeError) {
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
			this.GetDerived().(Server).handlePackage(pkg).Test()
		} else {
			this.packageBuffer = append(this.packageBuffer, data...)
		}
	}

	return
}

func (this commonServer) handlePacakge(data []byte) (err maybe.MaybeError) {
	err.Error(errors.New("calling abstract method:commonServer.handlePackage"))
	return
}
