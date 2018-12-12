package interfaces

import "unsafe"

type DenseTableElement interface {
	GetSize() int32
	Get(int64, unsafe.Pointer) bool
	Put(unsafe.Pointer, unsafe.Pointer) bool
	Erase(int64, unsafe.Pointer) bool
}
