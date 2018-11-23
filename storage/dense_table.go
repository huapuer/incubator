package storage

import (
	"fmt"
	"../common/maybe"
	"unsafe"
)

type sparseEntry struct {
	KeyTo  int64
    Offset int64
	Size   int64
	HashStride int32
}

type denseTable struct {
	blocksNum int64
	blockSize int64
	denseSize  int64
	sparseEntries []*sparseEntry
	data uintptr
	elementSize int32
	eraser []byte
	hashSteps int32
}

type MaybeDenseTable struct {
	maybe.MaybeError
	value *denseTable
}

func (this MaybeDenseTable) Value(value *denseTable) {
	this.Error(nil)
	this.value = value
}

func (this MaybeDenseTable) Right() *denseTable {
	this.Test()
	return this.value
}

type MaybeUintptr struct {
	maybe.MaybeError
	value uintptr
}

func (this MaybeUintptr) Value(value uintptr) {
	this.Error(nil)
	this.value = value
}

func (this MaybeUintptr) Right() uintptr {
	this.Test()
	return this.value
}

type mimicSlice struct {
	addr *unsafe.ArbitraryType
	len  int
	cap  int
}

func NewDenseTable(blocksNum int64, denseSize int64, sparseEntries []*sparseEntry, elementSize int32) (this MaybeDenseTable) {
	if blocksNum <=0 {
		this.Error(fmt.Errorf("illegal blocks num: %d", blocksNum))
		return
	}
	if elementSize <=0 {
		this.Error(fmt.Errorf("illegal element size: %d", elementSize))
		return
	}
	if denseSize < 0 {
		this.Error(fmt.Errorf("illegal dense size: %d", denseSize))
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

	totalSize := blocksNum * blockSize * int64(elementSize)
	bytes := make([]byte, totalSize, totalSize)

	this.Value(&denseTable{
		blocksNum:blocksNum,
		blockSize:blockSize,
		denseSize:denseSize,
		sparseEntries:sparseEntries,
		data: uintptr(unsafe.Pointer((*mimicSlice)(unsafe.Pointer(&bytes)).addr)),
		elementSize:elementSize,
		eraser: make([]byte, elementSize, elementSize),
	})
	return
}

func (this *denseTable) Get(block int64, key int64, hit func(element uintptr) bool) (ret MaybeUintptr) {
	if key < this.denseSize {
		ptr:= uintptr(block * this.blockSize + int64(this.data) + key*int64(this.elementSize))
		ret.Value(ptr)
		return
	}
	for _, entry := range this.sparseEntries {
		if key < entry.KeyTo {
			idx := key % entry.Size
			for i:=int32(0);i<this.hashSteps;i++ {
				ptr := uintptr(int64(this.data) + idx)
				if hit(ptr) {
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

func (this *denseTable) Put(block int64, key int64, val uintptr, hit func(element uintptr) bool) (err maybe.MaybeError) {
	if key < this.denseSize {
		ptr := uintptr(block * this.blockSize + int64(this.data) + key*int64(this.elementSize))
		src := *(*[]byte)(unsafe.Pointer(&mimicSlice{addr: (*unsafe.ArbitraryType)(val), len: int(this.elementSize)}))
		copy(*(*[]byte)(unsafe.Pointer(&mimicSlice{addr: (*unsafe.ArbitraryType)(ptr), len: int(this.elementSize)})), src)
		err.Error(nil)
		return
	}
	for _, entry := range this.sparseEntries {
		if key < entry.KeyTo {
			idx := key % entry.Size
			for i:=int32(0);i<this.hashSteps;i++ {
				ptr := uintptr(int64(this.data) + idx)
				if hit(ptr) {
					src := *(*[]byte)(unsafe.Pointer(&mimicSlice{addr: (*unsafe.ArbitraryType)(val), len: int(this.elementSize)}))
					copy(*(*[]byte)(unsafe.Pointer(&mimicSlice{addr: (*unsafe.ArbitraryType)(ptr), len: int(this.elementSize)})), src)
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

func (this *denseTable) Del(block int64, key int64, hit func(element uintptr) bool) (err maybe.MaybeError) {
	if key < this.denseSize {
		ptr:=uintptr(block * this.blockSize + int64(this.data) + key*int64(this.elementSize))
		copy(*(*[]byte)(unsafe.Pointer(&mimicSlice{addr: (*unsafe.ArbitraryType)(ptr), len: int(this.elementSize)})), this.eraser)
		err.Error(nil)
		return
	}
	for _, entry := range this.sparseEntries {
		if key < entry.KeyTo {
			idx := key % entry.Size
			for i:=int32(0);i<this.hashSteps;i++ {
				ptr := uintptr(int64(this.data) + idx)
				if hit(ptr) {
					copy(*(*[]byte)(unsafe.Pointer(&mimicSlice{addr: (*unsafe.ArbitraryType)(ptr), len: int(this.elementSize)})), this.eraser)
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
