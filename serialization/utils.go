package serialization

import "unsafe"

type mimicSlice struct {
	addr unsafe.Pointer
	len  int
	cap  int
}

type mimicIFace struct {
	tab  unsafe.Pointer
	data unsafe.Pointer
}

type mimicEFace struct {
	_type uintptr
	data  unsafe.Pointer
}

////go:noescape
func Ptr2IFace(iface unsafe.Pointer, ptr unsafe.Pointer) {
	(*mimicIFace)(iface).data = ptr
}

////go:noescape
func IFace2Ptr(iface unsafe.Pointer) unsafe.Pointer {
	return (*mimicIFace)(iface).data
}

////go:noescape
func Ptr2Bytes(ptr unsafe.Pointer, size int) []byte {
	return *(*[]byte)(unsafe.Pointer(&mimicSlice{addr: ptr, len: size, cap: size}))
}

////go:noescape
func Bytes2Ptr(src []byte) unsafe.Pointer {
	return unsafe.Pointer((*mimicSlice)(unsafe.Pointer(&src)).addr)
}

////go:noescape
func Move(dst unsafe.Pointer, src unsafe.Pointer, size int) {
	copy(*(*[]byte)(unsafe.Pointer(&mimicSlice{addr: dst, len: size, cap: size})),
		*(*[]byte)(unsafe.Pointer(&mimicSlice{addr: src, len: size, cap: size})))
}

////go:noescape
func MoveBytes(dst unsafe.Pointer, src []byte, size int) {
	copy(*(*[]byte)(unsafe.Pointer(&mimicSlice{addr: dst, len: size, cap: size})), src)
}

////go:noescape
func Eface2TypeInt(eface interface{}) int {
	return int((*mimicEFace)(unsafe.Pointer(&eface))._type)
}
