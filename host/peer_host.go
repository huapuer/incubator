package host

import (
	"../common/maybe"
	"../config"
	"../message"
	"../serialization"
	"errors"
	"fmt"
	"net"
)

const (
	defaultPeerHostClassName = "actor.peerHost"
)

func init() {
	RegisterHostPrototype(defaultPeerHostClassName, &peerHost{}).Test()
}

type peerHost struct {
	commonHost

	peers map[int64]net.Conn
}

func (this *peerHost) Receive(conn net.Conn, msg message.RemoteMessage) (err maybe.MaybeError) {
	m, ok := msg.(message.SeesionedMessage)
	if !ok {
		err.Error(fmt.Errorf("peer host receiving not sessioned message: %+v", msg))
		return
	}
	if m.IsToServer().Right() {
		if conn == nil {
			err.Error(errors.New("peer conn not set"))
			return
		}
		message.Route(m).Test()
		this.peers[m.GetSesseionId()] = conn
	} else {
		peer, ok := this.peers[m.GetSesseionId()]
		if !ok {
			err.Error(fmt.Errorf("peer not found: %d", m.GetSesseionId()))
			return
		}
		_, e := peer.Write(serialization.Marshal(m))
		if e != nil {
			delete(this.peers, m.GetSesseionId())
			err.Error(e)
			return
		}
	}
	delete(this.peers, m.GetSesseionId())
	err.Error(nil)
	return
}

func (this peerHost) New(attrs interface{}, cfg config.Config) config.IOC {
	ret := MaybeHost{}
	ret.Value(&peerHost{
		peers: make(map[int64]net.Conn),
	})
	return ret
}

func (this peerHost) GetJsonBytes() (ret maybe.MaybeBytes) {
	panic("not implemented")
}

func (this *peerHost) SetJsonField(data []byte) (err maybe.MaybeError) {
	panic("not implemented")
}

func (this peerHost) GetSize() int32 {
	panic("not implemented")
}
