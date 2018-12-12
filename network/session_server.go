package network

import (
	"github.com/incubator/common/maybe"
	"github.com/incubator/config"
	"github.com/incubator/interfaces"
	"github.com/incubator/message"
	"github.com/incubator/protocal"
	"github.com/incubator/serialization"
	"github.com/incubator/topo"
	"math/rand"
	"net"
)

const (
	sessionServerClassName = "network.sessionServer"
)

func init() {
	interfaces.RegisterServerPrototype(sessionServerClassName, &defaultServer{}).Test()
}

type sessionServer struct {
	commonServer
}

func (this sessionServer) New(attrs interface{}, cfg interfaces.Config) interfaces.IOC {
	ret := interfaces.MaybeServer{}

	network := config.GetAttrString(attrs, "Network", config.CheckStringNotEmpty).Right()
	protocalClass := config.GetAttrString(attrs, "Protocal", config.CheckStringNotEmpty).Right()

	s := &sessionServer{
		commonServer{
			network: network,
			p:       protocal.GetProtocalPrototype(protocalClass).Right(),
		},
	}
	s.Inherit(s)
	ret.Value(s)
	return ret
}

////go:noescape
func (this sessionServer) HandlePackage(data []byte, c net.Conn) (err maybe.MaybeError) {
	layerId := data[0]
	typ := data[1]

	l := interfaces.GetLayer(int32(layerId)).Right()
	msg := l.GetMessageCanonicalFromType(int32(typ)).Right()
	serialization.UnmarshalRemoteMessage(data, msg).Test()

	sessId := rand.Int63()

	l.GetTopo().(topo.SessionTopo).AddHost(sessId, c).Test()
	msg.(interfaces.RemoteMessage).SetHostId(sessId)
	message.Route(msg).Test()
	return
}
