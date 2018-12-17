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

	network       string
	port          int
	bufferSize    int
	packageBuffer []byte
	packageSize   int
	headerSize    int
	p             interfaces.Protocal
	derived       interfaces.Server
	handlerNum    int
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

	this.packageBuffer = make([]byte, 0, this.bufferSize)

	l, e := net.Listen(this.network, fmt.Sprintf("0.0.0.0:%d", this.port))
	if e != nil {
		err.Error(e)
		return
	}
	rand.Seed(time.Now().Unix())

	go func() {
		defer l.Close()
		for {
			select {
			case <-ctx.Done():
				err.Error(nil)
				return
			default:
				break
			}

			c, e := l.Accept()
			if e != nil {
				err.Error(e)
				return
			}
			for i := 0; i < this.handlerNum; i++ {
				go this.HandleConnection(ctx, c)
			}
		}
	}()

	err.Error(nil)
	return
}

func (this commonServer) HandleConnection(ctx context.Context, c net.Conn) {
	defer c.Close()
	reader := bufio.NewReader(c)
	buffer := make([]byte, this.bufferSize, this.bufferSize)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			break
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
func (this *commonServer) HandleData(data []byte, l int, c net.Conn) (err maybe.MaybeError) {
	if l == 0 {
		err.Error(errors.New("empty data"))
		return
	}
	this.packageBuffer = append(this.packageBuffer, data...)
	for {
		this.packageSize, this.headerSize = this.p.Parse(this.packageBuffer)
		if this.packageSize == protocal.PROTOCAL_PARSE_STATE_SHORT {
			break
		} else if this.packageSize == protocal.PROTOCAL_PARSE_STATE_ERROR {
			this.packageSize = protocal.PROTOCAL_PARSE_STATE_SHORT
			this.packageBuffer = make([]byte, this.bufferSize, 0)
			break
		}

		pkg := this.packageBuffer[this.headerSize:this.packageSize]
		this.packageBuffer = this.packageBuffer[this.packageSize:]
		this.derived.HandlePackage(this.p.Decode(pkg), c).Test()
	}

	err.Error(nil)
	return
}

func (this *commonServer) Inherit(cobj class.Class) {
	this.derived = cobj.(interfaces.Server)
}

func (this *commonServer) SetPort(port int) {
	this.port = port
}

func (this commonServer) GetProtocal() interfaces.Protocal {
	return this.p
}
