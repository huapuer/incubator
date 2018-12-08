package network

import (
	"../common/maybe"
	"../config"
	"../layer"
	"../message"
	"../protocal"
	"../serialization"
	"net"
)

const (
	defaultServerClassName = "server.defaultServer"
)

func init() {
	RegisterServerPrototype(defaultServerClassName, &defaultServer{}).Test()
}

type defaultServer struct {
	commonServer
}

func (this defaultServer) New(attrs interface{}, cfg config.Config) config.IOC {
	ret := MaybeServer{}

	network := config.GetAttrString(attrs, "Network", config.CheckStringNotEmpty).Right()
	port := config.GetAttrInt(attrs, "Port", config.CheckIntGT0).Right()
	handlerNum := config.GetAttrInt(attrs, "HandlerNum", config.CheckIntGT0).Right()
	protocalClass := config.GetAttrString(attrs, "Protocal", config.CheckStringNotEmpty).Right()

	s := &defaultServer{
		commonServer{
			network:    network,
			port:       port,
			handlerNum: handlerNum,
			p:          protocal.GetProtocalPrototype(protocalClass).Right(),
		},
	}
	s.Inherit(s)
	ret.Value(s)
	return ret
}

//go:noescape
func (this defaultServer) handlePackage(data []byte, c net.Conn) (err maybe.MaybeError) {
	layerId := data[0]
	typ := data[1]

	l := layer.GetLayer(int32(layerId)).Right()
	msg := l.GetMessageCanonicalFromType(int32(typ)).Right()
	serialization.UnmarshalRemoteMessage(data, msg).Test()

	message.Route(msg).Test()
	return
}
