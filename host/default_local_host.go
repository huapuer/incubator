package host

import (
	"../common/maybe"
	"../config"
	"../message"
	"github.com/incubator/serialization"
	"github.com/incubator/storage"
	"net"
	"sync/atomic"
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

func (this defaultLocalHost) GetStatePoint() *int64 {
	return &this.id
}

func (this defaultLocalHost) GetSize() int32 {
	return int32(unsafe.Sizeof(this))
}

func (this defaultLocalHost) Aquire(key int64, ptr unsafe.Pointer) bool {
	var h LocalHost
	h = &defaultLocalHost{}
	serialization.Ptr2IFace(&h, ptr)
	return atomic.CompareAndSwapInt64(h.GetStatePoint(), key, storage.DENSE_TABLE_ELEMENT_STATE_ONREAD)
}

func (this defaultLocalHost) Release(key int64, ptr unsafe.Pointer) bool {
	var h LocalHost
	h = &defaultLocalHost{}
	serialization.Ptr2IFace(&h, ptr)
	return atomic.CompareAndSwapInt64(h.GetStatePoint(), storage.DENSE_TABLE_ELEMENT_STATE_ONREAD, key)
}

func (this defaultLocalHost) Put(dst unsafe.Pointer, src unsafe.Pointer) bool {
	var h LocalHost
	h = &defaultLocalHost{}
	serialization.Ptr2IFace(&h, dst)
	if atomic.CompareAndSwapInt64(h.GetStatePoint(), storage.DENSE_TABLE_ELEMENT_STATE_EMPTY, storage.DENSE_TABLE_ELEMENT_STATE_ONWRITE) {
		var s LocalHost
		s = &defaultLocalHost{}
		serialization.Ptr2IFace(&s, src)
		id := s.GetId()
		s.SetId(-2)
		serialization.Move(dst, src, int(this.GetSize()))
		return atomic.CompareAndSwapInt64(h.GetStatePoint(), storage.DENSE_TABLE_ELEMENT_STATE_ONWRITE, id)
	}
	return false
}

func (this defaultLocalHost) Erase(key int64, ptr unsafe.Pointer) bool {
	var h LocalHost
	h = &defaultLocalHost{}
	serialization.Ptr2IFace(&h, ptr)
	return atomic.CompareAndSwapInt64(h.GetStatePoint(), key, storage.DENSE_TABLE_ELEMENT_STATE_EMPTY)
}
