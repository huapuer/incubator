package link

import (
	"github.com/incubator/serialization"
	"github.com/incubator/storage"
	"sync/atomic"
	"unsafe"
)

const (
	defaultLinkClassName = "link.defaultLink"
)

func init() {
	RegisterLinkPrototype(defaultLinkClassName, &defaultLink{}).Test()
}

type defaultLink struct {
	commonLink
}

func (this defaultLink) GetStatePoint() *int64 {
	return &this.toId
}

func (this defaultLink) GetSize() int32 {
	return int32(unsafe.Sizeof(this))
}

func (this defaultLink) Aquire(key int64, ptr unsafe.Pointer) bool {
	var h Link
	h = &defaultLink{}
	serialization.Ptr2IFace(&h, ptr)
	return atomic.CompareAndSwapInt64(h.GetStatePoint(), key, storage.DENSE_TABLE_ELEMENT_STATE_ONREAD)
}

func (this defaultLink) Release(key int64, ptr unsafe.Pointer) bool {
	var h Link
	h = &defaultLink{}
	serialization.Ptr2IFace(&h, ptr)
	return atomic.CompareAndSwapInt64(h.GetStatePoint(), storage.DENSE_TABLE_ELEMENT_STATE_ONREAD, key)
}

func (this defaultLink) Put(dst unsafe.Pointer, src unsafe.Pointer) bool {
	var h Link
	h = &defaultLink{}
	serialization.Ptr2IFace(&h, dst)
	if atomic.CompareAndSwapInt64(h.GetStatePoint(), storage.DENSE_TABLE_ELEMENT_STATE_EMPTY, storage.DENSE_TABLE_ELEMENT_STATE_ONWRITE) {
		var s Link
		s = &defaultLink{}
		serialization.Ptr2IFace(&s, src)
		id := s.GetToId()
		s.SetToId(-2)
		serialization.Move(dst, src, int(this.GetSize()))
		return atomic.CompareAndSwapInt64(h.GetStatePoint(), storage.DENSE_TABLE_ELEMENT_STATE_ONWRITE, id)
	}
	return false
}

func (this defaultLink) Erase(key int64, ptr unsafe.Pointer) bool {
	var h Link
	h = &defaultLink{}
	serialization.Ptr2IFace(&h, ptr)
	return atomic.CompareAndSwapInt64(h.GetStatePoint(), key, storage.DENSE_TABLE_ELEMENT_STATE_EMPTY)
}
