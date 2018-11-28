package network

import (
	"../common/maybe"
	"../layer"
	"../message"
	"github.com/incubator/serialization"
	"net"
)

type defaultServer struct {
	commonServer
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
