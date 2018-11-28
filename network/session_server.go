package network

import (
	"../common/maybe"
	"../layer"
	"github.com/incubator/message"
	"github.com/incubator/serialization"
	"github.com/incubator/topo"
	"math/rand"
	"net"
)

type sessionServer struct {
	commonServer
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
	msg.(message.SeesionedMessage).SetSessionId(sessId)
	message.Route(msg).Test()
	return
}
