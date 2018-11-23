package host

import (
	"../common/maybe"
	"../config"
	"../message"
	"github.com/incubator/serialization"
	"net"
	"unsafe"
)

const (
	defaultLocalHostClassName = "actor.defaultLocalHost"
)

func init() {
	RegisterHostPrototype(defaultLocalHostClassName, &defaultLocalHost{}).Test()
}

type defaultLocalHost struct {
	commonHost
}

func (this *defaultLocalHost) Receive(conn net.Conn, msg message.RemoteMessage) (err maybe.MaybeError) {
	message.Route(msg).Test()
	return
}

func (this defaultLocalHost) New(attrs interface{}, cfg config.Config) config.IOC {
	ret := MaybeHost{}
	//TODO: real logic
	ret.Value(&defaultLocalHost{
		commonHost{
			valid: true,
		},
	})
	return ret
}

func (this defaultLocalHost) IsHit(key int64, ptr unsafe.Pointer) bool {
	var h LocalHost
	h = &defaultLocalHost{}
	serialization.Ptr2IFace(&h, ptr)
	return h.GetId().Right() == key
}

func (this defaultLocalHost) IsEmpty(ptr unsafe.Pointer) bool {
	var h LocalHost
	h = &defaultLocalHost{}
	serialization.Ptr2IFace(&h, ptr)
	return h.GetId().Right() >= 0
}

func (this defaultLocalHost) Erase(ptr unsafe.Pointer) {
	var h LocalHost
	h = &defaultLocalHost{}
	serialization.Ptr2IFace(&h, ptr)
	h.SetId(-1)
}

func (this defaultLocalHost) GetSize() int32 {
	return int32(unsafe.Sizeof(this))
}
