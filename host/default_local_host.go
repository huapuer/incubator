package host

import (
	"github.com/incubator/common/maybe"
	"github.com/incubator/interfaces"
	"github.com/incubator/message"
	"github.com/incubator/serialization"
	"github.com/incubator/storage"
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

func (this defaultLocalHost) Receive(msg interfaces.RemoteMessage) (err maybe.MaybeError) {
	message.Route(msg).Test()
	return
}

func (this defaultLocalHost) New(attrs interface{}, cfg interfaces.Config) interfaces.IOC {
	ret := interfaces.MaybeHost{}
	//TODO: real logic
	ret.Value(&defaultLocalHost{})
	return ret
}

func (this defaultLocalHost) GetSize() int32 {
	return int32(unsafe.Sizeof(this))
}

func (this defaultLocalHost) Get(key int64, ptr unsafe.Pointer) bool {
	var h interfaces.Host
	h = &defaultLocalHost{}
	serialization.Ptr2IFace(&h, ptr)
	return h.GetId() == key
}

func (this defaultLocalHost) Put(dst unsafe.Pointer, src unsafe.Pointer) bool {
	var h interfaces.Host
	h = &defaultLocalHost{}
	serialization.Ptr2IFace(&h, dst)
	if h.GetId() == storage.DENSE_TABLE_ELEMENT_STATE_EMPTY {
		serialization.Move(dst, src, int(this.GetSize()))
		return true
	}
	return false
}

func (this defaultLocalHost) Erase(key int64, ptr unsafe.Pointer) bool {
	var h interfaces.Host
	h = &defaultLocalHost{}
	serialization.Ptr2IFace(&h, ptr)
	if h.GetId() == key {
		h.SetId(storage.DENSE_TABLE_ELEMENT_STATE_EMPTY)
	}
	return true
}

func (this defaultLocalHost) IsHealth() bool {
	panic("not implemented")
}
