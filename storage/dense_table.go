package storage

import (
	"../common/maybe"
	"../serialization"
	"errors"
	"fmt"
	"sync/atomic"
	"unsafe"
)

const (
	DENSE_TABLE_ELEMENT_STATE_EMPTY = -1
)

type SparseEntry struct {
	KeyTo      int64
	Offset     int64
	Size       int64
	HashStride int32
}

type DenseTable struct {
	elemLen       int64
	blockLen      int64
	blockSize     int64
	denseSize     int64
	sparseEntries []*SparseEntry
	data          uintptr
	size          int64
	elementSize   int32
	hashDepth     int32
	elementCanon  DenseTableElement
}

type MaybeDenseTable struct {
	maybe.MaybeError
	value DenseTable
}

func (this MaybeDenseTable) Value(value DenseTable) {
	this.Error(nil)
	this.value = value
}

func (this MaybeDenseTable) Right() DenseTable {
	this.Test()
	return this.value
}

type DenseTableElement interface {
	GetSize() int32
	Get(int64, unsafe.Pointer) bool
	Put(unsafe.Pointer, unsafe.Pointer) bool
	Erase(int64, unsafe.Pointer) bool
}

type MaybePointer struct {
	maybe.MaybeError
	value unsafe.Pointer
}

func (this MaybePointer) Value(value unsafe.Pointer) {
	this.Error(nil)
	this.value = value
}

func (this MaybePointer) Right() unsafe.Pointer {
	this.Test()
	return this.value
}

func NewDenseTable(elementCanon DenseTableElement,
	blocksNum int64,
	denseSize int64,
	sparseEntries []*SparseEntry,
	elementSize int32,
	hashDepth int32,
	data []byte) (this MaybeDenseTable) {
	if elementCanon == nil {
		this.Error(errors.New("element canon is nil"))
		return
	}
	if blocksNum <= 0 {
		this.Error(fmt.Errorf("illegal blocks num: %d", blocksNum))
		return
	}
	if elementSize <= 0 {
		this.Error(fmt.Errorf("illegal element size: %d", elementSize))
		return
	}
	if denseSize < 0 {
		this.Error(fmt.Errorf("illegal dense size: %d", denseSize))
		return
	}
	if hashDepth <= 0 {
		this.Error(fmt.Errorf("illegal hash steps: %d", hashDepth))
		return
	}
	blockSize := denseSize
	currentKeyOffset := int64(0)
	for _, entry := range sparseEntries {
		if entry.KeyTo <= currentKeyOffset {
			this.Error(fmt.Errorf("key-to not in asec order: %d", entry.KeyTo))
			return
		}
		if entry.Size <= 0 {
			this.Error(fmt.Errorf("illegal sparse entry size: %d", entry.Size))
			return
		}
		if entry.KeyTo < entry.Size {
			this.Error(fmt.Errorf("key-to less than size: %d<%d", entry.KeyTo, entry.Size))
			return
		}
		entry.Offset = blockSize
		blockSize += entry.Size
	}

	size := int64(0)
	if data == nil {
		size = blocksNum * blockSize * int64(elementSize)
		data = make([]byte, size, size)
	} else {
		size = int64(len(data))
	}

	this.Value(DenseTable{
		blockLen:      blocksNum,
		blockSize:     blockSize,
		denseSize:     denseSize,
		sparseEntries: sparseEntries,
		data:          uintptr(serialization.Bytes2Ptr(data)),
		elementSize:   elementSize,
		hashDepth:     hashDepth,
		size:          size,
	})
	return
}

func (this DenseTable) GetBytes() []byte {
	return serialization.Ptr2Bytes(unsafe.Pointer(this.data), int(this.size))
}

func (this *DenseTable) Get(block int64, key int64) (ret MaybePointer) {
	if key < this.denseSize {
		ptr := unsafe.Pointer(uintptr(block*this.blockSize + int64(this.data) + key*int64(this.elementSize)))
		ret.Value(ptr)
		return
	}
	for _, entry := range this.sparseEntries {
		if key < entry.KeyTo {
			idx := key % entry.Size
			for i := int32(0); i < this.hashDepth; i++ {
				ptr := unsafe.Pointer(uintptr(int64(this.data) + idx))
				if this.elementCanon.Get(key, ptr) {
					ret.Value(ptr)
					return
				}
				idx += int64(entry.HashStride)
				if idx >= entry.Size {
					idx %= entry.Size
				}
			}
		}
	}
	ret.Error(fmt.Errorf("key not found: %d", key))
	return
}

//go:noescape
func (this *DenseTable) Put(block int64, key int64, val unsafe.Pointer) bool {
	if key < this.denseSize {
		ptr := unsafe.Pointer(uintptr(block*this.blockSize + int64(this.data) + key*int64(this.elementSize)))
		serialization.Move(ptr, val, int(this.elementSize))
		atomic.AddInt64(&this.elemLen, 1)
		return true
	}
	for _, entry := range this.sparseEntries {
		if key < entry.KeyTo {
			idx := key % entry.Size
			for i := int32(0); i < this.hashDepth; i++ {
				ptr := unsafe.Pointer(uintptr(int64(this.data) + idx))
				if this.elementCanon.Put(ptr, val) {
					serialization.Move(ptr, val, int(this.elementSize))
					atomic.AddInt64(&this.elemLen, 1)
					return true
				}
				idx += int64(entry.HashStride)
				if idx >= entry.Size {
					idx %= entry.Size
				}
			}
		}
	}
	return false
}

func (this *DenseTable) Del(block int64, key int64) bool {
	if key < this.denseSize {
		ptr := unsafe.Pointer(uintptr(block*this.blockSize + int64(this.data) + key*int64(this.elementSize)))
		this.elementCanon.Erase(key, ptr)
		atomic.AddInt64(&this.elemLen, -1)
		return true
	}
	for _, entry := range this.sparseEntries {
		if key < entry.KeyTo {
			idx := key % entry.Size
			for i := int32(0); i < this.hashDepth; i++ {
				ptr := unsafe.Pointer(uintptr(int64(this.data) + idx))
				if this.elementCanon.Erase(key, ptr) {
					atomic.AddInt64(&this.elemLen, -1)
					return true
				}
				idx += int64(entry.HashStride)
				if idx >= entry.Size {
					idx %= entry.Size
				}
			}
		}
	}
	return false
}

func (this *DenseTable) TraverseBlock(block int64, callback func(ptr unsafe.Pointer) bool) {
	base := block * this.blockSize * int64(this.elementSize)
	for i := int64(0); i < this.blockSize; i++ {
		ptr := unsafe.Pointer(uintptr(int64(this.data) + base + i*int64(this.elementSize)))
		if !callback(ptr) {
			return
		}
	}
	return
}

func (this DenseTable) BlockLen() int64 {
	return this.blockLen
}

func (this DenseTable) ElemLen() int64 {
	return this.elemLen
}
