package storage

import (
	"../common/maybe"
	"errors"
	"fmt"
	"../serialization"
	"sync/atomic"
	"unsafe"
)

const (
	DENSE_TABLE_ELEMENT_STATE_ONREAD  = -2
	DENSE_TABLE_ELEMENT_STATE_ONWRITE = -1
	DENSE_TABLE_ELEMENT_STATE_EMPTY   = 0
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
	GetStatePoint() *int64
	GetSize() int32
	Aquire(int64, unsafe.Pointer) bool
	Release(int64, unsafe.Pointer) bool
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

	if data == nil {
		totalSize := blocksNum * blockSize * int64(elementSize)
		data = make([]byte, totalSize, totalSize)
	}

	this.Value(DenseTable{
		blockLen:      blocksNum,
		blockSize:     blockSize,
		denseSize:     denseSize,
		sparseEntries: sparseEntries,
		data:          uintptr(serialization.Bytes2Ptr(data)),
		elementSize:   elementSize,
		hashDepth:     hashDepth,
	})
	return
}

func (this *DenseTable) Aquire(block int64, key int64) (ret MaybePointer) {
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
				if this.elementCanon.Aquire(key, ptr) {
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

func (this *DenseTable) Release(key int64, ptr unsafe.Pointer) bool {
	if key < this.denseSize {
		return true
	}
	return this.elementCanon.Release(key, ptr)
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

func (this DenseTable) BlockLen() int64 {
	return this.blockLen
}

func (this DenseTable) ElemLen() int64 {
	return this.elemLen
}
