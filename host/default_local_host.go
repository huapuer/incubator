package host

import (
	"../common/maybe"
	"../config"
	"../message"
	"../serialization"
	"../storage"
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

func (this defaultLocalHost) Receive(msg message.RemoteMessage) (err maybe.MaybeError) {
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

func (this defaultLocalHost) GetSize() int32 {
	return int32(unsafe.Sizeof(this))
}

func (this defaultLocalHost) Get(key int64, ptr unsafe.Pointer) bool {
	var h Host
	h = &defaultLocalHost{}
	serialization.Ptr2IFace(&h, ptr)
	return h.GetId() == key
}

func (this defaultLocalHost) Put(dst unsafe.Pointer, src unsafe.Pointer) bool {
	var h Host
	h = &defaultLocalHost{}
	serialization.Ptr2IFace(&h, dst)
	if h.GetId() == storage.DENSE_TABLE_ELEMENT_STATE_EMPTY {
		serialization.Move(dst, src, int(this.GetSize()))
		return true
	}
	return false
}

func (this defaultLocalHost) Erase(key int64, ptr unsafe.Pointer) bool {
	var h Host
	h = &defaultLocalHost{}
	serialization.Ptr2IFace(&h, ptr)
	if h.GetId() == key {
		h.SetId(storage.DENSE_TABLE_ELEMENT_STATE_EMPTY)
	}
	return true
}
