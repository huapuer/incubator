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
	defaultLinkHostClassName = "host.defaultLinkHost"
)

func init() {
	RegisterHostPrototype(defaultLinkHostClassName, &defaultLinkHost{}).Test()
}

type defaultLinkHost struct {
	commonHost
	commonLinkHost
}

func (this *defaultLinkHost) Receive(msg message.RemoteMessage) (err maybe.MaybeError) {
	message.Route(msg).Test()
	return
}

func (this defaultLinkHost) New(attrs interface{}, cfg config.Config) config.IOC {
	ret := MaybeHost{}
	//TODO: real logic
	ret.Value(&defaultLinkHost{
		commonHost:     commonHost{},
		commonLinkHost: commonLinkHost{},
	})
	return ret
}

func (this defaultLinkHost) GetSize() int32 {
	return int32(unsafe.Sizeof(this))
}

func (this defaultLinkHost) Get(key int64, ptr unsafe.Pointer) bool {
	var h Host
	h = &defaultLinkHost{}
	serialization.Ptr2IFace(&h, ptr)
	return h.GetId() == key
}

func (this defaultLinkHost) Put(dst unsafe.Pointer, src unsafe.Pointer) bool {
	var h Host
	h = &defaultLinkHost{}
	serialization.Ptr2IFace(&h, dst)
	if h.GetId() == storage.DENSE_TABLE_ELEMENT_STATE_EMPTY {
		serialization.Move(dst, src, int(this.GetSize()))
		return true
	}
	return false
}

func (this defaultLinkHost) Erase(key int64, ptr unsafe.Pointer) bool {
	var h Host
	h = &defaultLinkHost{}
	serialization.Ptr2IFace(&h, ptr)
	if h.GetId() == key {
		h.SetId(storage.DENSE_TABLE_ELEMENT_STATE_EMPTY)
	}
	return true
}

func (this defaultLinkHost) IsHealth() bool {
	panic("not implemented")
}
