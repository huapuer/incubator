package storage

import (
	"errors"
	"fmt"
	"../common/maybe"
)

type TableConfig struct {
	denseSize  int64
	sparseSize int64
}

type linkedUintPtr struct {
	this uintptr
	next *linkedUintPtr
}

type denseTable struct {
	denseTable  []uintptr
	sparseTable []*linkedUintPtr
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

type MaybeLinkedUintptr struct {
	maybe.MaybeError
	value *linkedUintPtr
}

func (this MaybeLinkedUintptr) Value(value *linkedUintPtr) {
	this.Error(nil)
	this.value = value
}

func (this MaybeLinkedUintptr) Right() *linkedUintPtr {
	this.Test()
	return this.value
}

func NewDenseTable(config *TableConfig) (this MaybeDenseTable) {
	if config == nil {
		this.Error(errors.New("empty config"))
		return
	}
	if config.denseSize <= 0 || config.sparseSize <= 0 {
		this.Error(fmt.Errorf("illegal config: denseSize=%d, sparseSize=%d", config.denseSize, config.sparseSize))
		return
	}
	this.Value(
		&denseTable{
			denseTable:  make([]uintptr, config.denseSize, config.denseSize),
			sparseTable: make([]*linkedUintPtr, config.sparseSize, config.sparseSize),
		})
	return
}

func (this *denseTable) GetDense(key int64) (ret MaybeUintptr) {
	if key < 0 || key > int64(len(this.denseTable)) {
		ret.Error(fmt.Errorf("illegal dense key:%d", key))
		return
	}
	if this.denseTable[key] != 0 {
		ret.Value(this.denseTable[key])
	}
	ret.Error(fmt.Errorf("dense key not found: %d", key))
	return
}

func (this *denseTable) PutDense(key int64, value uintptr) (ret maybe.MaybeError) {
	if key < 0 || key > int64(len(this.denseTable)) {
		ret.Error(fmt.Errorf("illegal dense key:%d", key))
		return
	}
	this.denseTable[key] = value
	return
}

func (this *denseTable) DeleteDense(key int64) (ret maybe.MaybeError) {
	if key < 0 || key > int64(len(this.denseTable)) {
		ret.Error(fmt.Errorf("illegal dense key:%d", key))
		return
	}
	this.denseTable[key] = 0
	return
}

func (this *denseTable) GetSparse(key int64, compare func(uintptr) bool) (ret MaybeLinkedUintptr) {
	if key < 0 {
		ret.Error(fmt.Errorf("illegal sparse key:%d", key))
		return
	}
	first := this.sparseTable[key%int64(len(this.sparseTable))]
	for {
		if first == nil {
			ret.Value(first)
			return
		}
		if compare(first.this) {
			ret.Value(first)
			return
		}
		first = first.next
	}
	return
}

func (this *denseTable) PutSparse(key int64, value *linkedUintPtr) (ret maybe.MaybeError) {
	if key < 0 {
		ret.Error(fmt.Errorf("illegal sparse key:%d", key))
		return
	}
	if value == nil {
		ret.Error(fmt.Errorf("value is nil: %s", key))
		return
	}
	value.next = this.sparseTable[key%int64(len(this.sparseTable))]
	this.sparseTable[key%int64(len(this.sparseTable))] = value
	return
}

func (this *denseTable) DeleteSparse(key int64, compare func(uintptr) bool) (ret maybe.MaybeError) {
	if key < 0 {
		ret.Error(fmt.Errorf("illegal sparse key:%d", key))
		return
	}
	first := this.sparseTable[key%int64(len(this.sparseTable))]
	if first == nil {
		ret.Error(fmt.Errorf("sparse key not found:%d", key))
		return
	}
	if compare(first.this) {
		this.sparseTable[key%int64(len(this.sparseTable))] = first.next
		return
	}
	for {
		if first.next == nil {
			ret.Error(fmt.Errorf("sparse key not found:%d", key))
			return
		}
		if compare(first.next.this) {
			first.next = first.next.next
			return
		}
		first = first.next
	}
	return
}
