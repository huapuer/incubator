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
	defaultLinkHostClassName = "host.defaultLinkHost"
)

func init() {
	RegisterHostPrototype(defaultLinkHostClassName, &defaultLinkHost{}).Test()
}

type defaultLinkHost struct {
	commonHost
	commonLinkHost
}

func (this *defaultLinkHost) Receive(msg interfaces.RemoteMessage) (err maybe.MaybeError) {
	message.Route(msg).Test()
	return
}

func (this defaultLinkHost) New(attrs interface{}, cfg interfaces.Config) interfaces.IOC {
	ret := interfaces.MaybeHost{}
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
	var h interfaces.Host
	h = &defaultLinkHost{}
	serialization.Ptr2IFace(&h, ptr)
	return h.GetId() == key
}

func (this defaultLinkHost) Put(dst unsafe.Pointer, src unsafe.Pointer) bool {
	var h interfaces.Host
	h = &defaultLinkHost{}
	serialization.Ptr2IFace(&h, dst)
	if h.GetId() == storage.DENSE_TABLE_ELEMENT_STATE_EMPTY {
		serialization.Move(dst, src, int(this.GetSize()))
		return true
	}
	return false
}

func (this defaultLinkHost) Erase(key int64, ptr unsafe.Pointer) bool {
	var h interfaces.Host
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
