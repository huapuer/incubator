package topo

import (
	"github.com/incubator/common/maybe"
	"github.com/incubator/interfaces"
	"net"
)

type SessionTopo interface {
	interfaces.Topo

	AddHost(int64, net.Conn) maybe.MaybeError
}

type commonTopo struct {
	layer int32
}

func (this commonTopo) GetLayer() int32 {
	return this.layer
}

func (this *commonTopo) SetLayer(layer int32) {
	this.layer = layer
}
