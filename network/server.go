package network

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/incubator/common/class"
	"github.com/incubator/common/maybe"
	"github.com/incubator/global"
	"github.com/incubator/interfaces"
	"github.com/incubator/protocal"
	"math/rand"
	"net"
	"time"
)

type commonServer struct {
	class.Class

	network        string
	port           int
	readBufferSize int
	packageBuffer  []byte
	packageSize    int
	headerSize     int
	p              interfaces.Protocal
	derived        interfaces.Server
	handlerNum     int
}

func (this commonServer) Start(ctx context.Context) (err maybe.MaybeError) {
	global.AddListenedPort(this.port).Test()

	if this.network == "" {
		err.Error(errors.New("not network provided"))
		return
	}
	if this.port <= 0 {
		err.Error(fmt.Errorf("illegal listen port: %d", this.port))
		return
	}
	if this.handlerNum <= 0 {
		err.Error(fmt.Errorf("illegal handler num: %d", this.handlerNum))
		return
	}

	this.packageBuffer = make([]byte, this.readBufferSize, 0)

	l, e := net.Listen(this.network, fmt.Sprint("%s:%d", "0.0.0.0", this.port))
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
			for i := 0; i < this.handlerNum; i++ {
				go this.HandleConnection(ctx, c)
			}
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
		this.HandleData(buffer, len, c).Test()
	}
	return
}

////go:noescape
func (this commonServer) HandlePacakge(data []byte, c net.Conn) (err maybe.MaybeError) {
	this.derived.HandlePackage(data, c)
	return
}

////go:noescape
func (this commonServer) HandleData(data []byte, l int, c net.Conn) (err maybe.MaybeError) {
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
			this.derived.HandlePackage(this.p.Decode(pkg), c).Test()
			this.packageSize = protocal.PROTOCAL_PARSE_STATE_SHORT
		}
	} else if this.packageSize == protocal.PROTOCAL_PARSE_STATE_ERROR {
		this.packageSize = protocal.PROTOCAL_PARSE_STATE_SHORT
		this.packageBuffer = make([]byte, this.readBufferSize, 0)
	}

	return
}

func (this *commonServer) Inherit(cobj class.Class) {
	this.derived = cobj.(interfaces.Server)
}

func (this *commonServer) SetPort(port int) {
	this.port = port
}
