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
	s := &defaultServer{
		commonServer{
			network: cfg.Server.Network,
			port:    cfg.Server.Port,
			p:       protocal.GetProtocalPrototype(cfg.Server.Protocal).Right(),
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
