package storage

import (
	"../common/maybe"
	"errors"
	"fmt"
	"github.com/incubator/serialization"
	"sync/atomic"
	"unsafe"
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
	GetSize() int32
	IsHit(int64, unsafe.Pointer) bool
	IsEmpty(unsafe.Pointer) bool
	Erase(unsafe.Pointer)
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
				if this.elementCanon.IsHit(key, ptr) {
					ret.Value(ptr)
					return
				}
				idx += int64(entry.HashStride)
			}
		}
	}
	ret.Error(fmt.Errorf("key not found: %d", key))
	return
}

//go:noescape
func (this *DenseTable) Put(block int64, key int64, val unsafe.Pointer) (err maybe.MaybeError) {
	if key < this.denseSize {
		ptr := unsafe.Pointer(uintptr(block*this.blockSize + int64(this.data) + key*int64(this.elementSize)))
		serialization.Move(ptr, val, int(this.elementSize))
		atomic.AddInt64(&this.elemLen, 1)
		err.Error(nil)
		return
	}
	for _, entry := range this.sparseEntries {
		if key < entry.KeyTo {
			idx := key % entry.Size
			for i := int32(0); i < this.hashDepth; i++ {
				ptr := unsafe.Pointer(uintptr(int64(this.data) + idx))
				if this.elementCanon.IsEmpty(ptr) {
					serialization.Move(ptr, val, int(this.elementSize))
					atomic.AddInt64(&this.elemLen, 1)
					err.Error(nil)
					return
				}
				idx += int64(entry.HashStride)
			}
		}
	}
	err.Error(fmt.Errorf("putting key failed: %d", key))
	return
}

func (this *DenseTable) Del(block int64, key int64) (err maybe.MaybeError) {
	if key < this.denseSize {
		ptr := unsafe.Pointer(uintptr(block*this.blockSize + int64(this.data) + key*int64(this.elementSize)))
		this.elementCanon.Erase(ptr)
		atomic.AddInt64(&this.elemLen, -1)
		err.Error(nil)
		return
	}
	for _, entry := range this.sparseEntries {
		if key < entry.KeyTo {
			idx := key % entry.Size
			for i := int32(0); i < this.hashDepth; i++ {
				ptr := unsafe.Pointer(uintptr(int64(this.data) + idx))
				if this.elementCanon.IsHit(key, ptr) {
					this.elementCanon.Erase(ptr)
					atomic.AddInt64(&this.elemLen, -1)
					err.Error(nil)
					return
				}
				idx += int64(entry.HashStride)
			}
		}
	}
	err.Error(fmt.Errorf("key not found: %d", key))
	return
}

func (this DenseTable) BlockLen() int64 {
	return this.blockLen
}

func (this DenseTable) ElemLen() int64 {
	return this.elemLen
}
