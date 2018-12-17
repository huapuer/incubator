package network

import (
	"github.com/incubator/common/maybe"
	"github.com/incubator/config"
	"github.com/incubator/interfaces"
	"github.com/incubator/message"
	"github.com/incubator/protocal"
	"github.com/incubator/serialization"
	"net"
)

const (
	defaultServerClassName = "network.defaultServer"
)

func init() {
	interfaces.RegisterServerPrototype(defaultServerClassName, &defaultServer{}).Test()
}

type defaultServer struct {
	commonServer
}

func (this defaultServer) New(attrs interface{}, cfg interfaces.Config) interfaces.IOC {
	ret := interfaces.MaybeServer{}

	network := config.GetAttrString(attrs, "Network", config.CheckStringNotEmpty).Right()
	handlerNum := config.GetAttrInt(attrs, "HandlerNum", config.CheckIntGT0).Right()
	protocalClass := config.GetAttrString(attrs, "Protocal", config.CheckStringNotEmpty).Right()
	bufferSize := config.GetAttrInt(attrs, "BufferSize", config.CheckIntGT0).Right()

	s := &defaultServer{
		commonServer{
			network:    network,
			handlerNum: handlerNum,
			p:          protocal.GetProtocalPrototype(protocalClass).Right(),
			bufferSize: bufferSize,
		},
	}
	s.Inherit(s)
	ret.Value(s)
	return ret
}

////go:noescape
func (this defaultServer) HandlePackage(data []byte, c net.Conn) (err maybe.MaybeError) {
	layerId := data[0]
	typ := data[1]

	l := interfaces.GetLayer(int32(layerId)).Right()
	msg := l.GetMessageCanonicalFromType(int32(typ)).Right()
	serialization.UnmarshalRemoteMessage(data, &msg).Test()

	message.Route(msg).Test()

	err.Error(nil)
	return
}
