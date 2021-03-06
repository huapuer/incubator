package host

import (
	"github.com/incubator/common/maybe"
	"github.com/incubator/interfaces"
	"github.com/incubator/serialization"
	"net"
	"unsafe"
)

const (
	defaultSessionHostClassName = "host.defaultSessionHost"
)

func init() {
	interfaces.RegisterHostPrototype(defaultSessionHostClassName, &defaultSessionHost{}).Test()
}

type defaultSessionHost struct {
	commonHost

	peer net.Conn
}

func (this *defaultSessionHost) Receive(msg interfaces.RemoteMessage) (err maybe.MaybeError) {
	_, e := this.peer.Write(serialization.Marshal(msg))
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

func (this defaultSessionHost) New(attrs interface{}, cfg interfaces.Config) interfaces.IOC {
	ret := interfaces.MaybeHost{}
	ret.Value(&defaultSessionHost{})
	return ret
}

func (this *defaultSessionHost) Replicate() (ret interfaces.MaybeHost) {
	new := *this
	ret.Value(&new)
	return
}

func (this defaultSessionHost) GetSize() int32 {
	panic("not implemented")
}

func (this defaultSessionHost) Get(key int64, ptr unsafe.Pointer) bool {
	panic("not implemented")
}

func (this defaultSessionHost) Put(dst unsafe.Pointer, src unsafe.Pointer) bool {
	panic("not implemented")
}

func (this defaultSessionHost) Erase(key int64, ptr unsafe.Pointer) bool {
	panic("not implemented")
}

func (this defaultSessionHost) IsHealth() bool {
	panic("not implemented")
}
