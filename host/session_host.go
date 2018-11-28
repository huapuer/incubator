package host

import (
	"../common/maybe"
	"../config"
	"../message"
	"../serialization"
	"fmt"
	"net"
)

const (
	defaultSessionHostClassName = "actor.defaultSessionHost"
)

func init() {
	RegisterHostPrototype(defaultSessionHostClassName, &defaultSessionHost{}).Test()
}

type defaultSessionHost struct {
	commonHost

	peer net.Conn
}

func (this *defaultSessionHost) Receive(msg message.RemoteMessage) (err maybe.MaybeError) {
	m, ok := msg.(message.SeesionedMessage)
	if !ok {
		err.Error(fmt.Errorf("peer host receiving not sessioned message: %+v", msg))
		return
	}

	_, e := this.peer.Write(serialization.Marshal(m))
	if e != nil {
		err.Error(e)
		return
	}

	this.peer.Close()

	err.Error(nil)
	return
}

func (this *defaultSessionHost) SetPeer(conn net.Conn) {
	this.peer = conn
	return
}

func (this defaultSessionHost) New(attrs interface{}, cfg config.Config) config.IOC {
	ret := MaybeHost{}
	ret.Value(&defaultSessionHost{})
	return ret
}

func (this *defaultSessionHost) Replicate() (ret MaybeHost) {
	new := *this
	ret.Value(&new)
	return
}
