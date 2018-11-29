package network

import (
	"../common/maybe"
	"../layer"
	"../message"
	"../serialization"
	"../topo"
	"math/rand"
	"net"
	"incubator/config"
	"incubator/protocal"
)

const (
	sessionServerClassName = "server.sessionServer"
)

func init() {
	RegisterServerPrototype(sessionServerClassName, &defaultServer{}).Test()
}

type sessionServer struct {
	commonServer
}

func (this sessionServer) New(attrs interface{}, cfg config.Config) config.IOC {
	ret := MaybeServer{}
	s := &sessionServer{
		commonServer{
			network:cfg.Server.Network,
			address:cfg.Server.Address,
			p: protocal.GetProtocalPrototype(cfg.Server.Protocal).Right(),
		},
	}
	s.Inherit(s)
	ret.Value(s)
	return ret
}

//go:noescape
func (this sessionServer) handlePackage(data []byte, c net.Conn) (err maybe.MaybeError) {
	layerId := data[0]
	typ := data[1]

	l := layer.GetLayer(int32(layerId)).Right()
	msg := l.GetMessageCanonicalFromType(int32(typ)).Right()
	serialization.UnmarshalRemoteMessage(data, msg).Test()

	sessId := rand.Int63()

	l.GetTopo().(topo.SessionTopo).AddHost(sessId, c).Test()
	msg.(message.RemoteMessage).SetHostId(sessId)
	message.Route(msg).Test()
	return
}
