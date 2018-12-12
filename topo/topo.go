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
