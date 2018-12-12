// +build !debug

package maybe

import (
	//log "github.com/sirupsen/logrus"
	"log"
)

type NilValueError struct {
	msg string
}

func (this NilValueError) Error() string {
	return this.msg
}

func TryCatch(try func(), catch func(error)) {
	defer func() {
		if err := recover(); err != nil {
			log.Fatalf("err=%+v", err)
			if catch != nil {
				catch(err.(error))
			}
		}
	}()
	try()
}

type Maybe interface {
	Error(error)
	Test()
}

type MaybeError struct {
	nonNil bool
	err    error
}

func (this *MaybeError) Error(err error) {
	this.err = err
	if err == nil {
		this.nonNil = true
	}
}

func (this MaybeError) Test() {
	if this.err != nil {
		panic(this.err)
	}

	if this.nonNil == false {
		panic(NilValueError{"Value not set."})
	}
}

// MaybeBool
type MaybeBool struct {
	MaybeError
	value bool
}

func (this *MaybeBool) Value(value bool) {
	this.Error(nil)
	this.value = value
}

func (this MaybeBool) Right() bool {
	this.Test()
	return this.value
}

// MaybeInt
type MaybeInt struct {
	MaybeError
	value int
}

func (this *MaybeInt) Value(value int) {
	this.Error(nil)
	this.value = value
}

func (this MaybeInt) Right() int {
	this.Test()
	return this.value
}

// MaybeInt32
type MaybeInt32 struct {
	MaybeError
	value int32
}

func (this *MaybeInt32) Value(value int32) {
	this.Error(nil)
	this.value = value
}

func (this MaybeInt32) Right() int32 {
	this.Test()
	return this.value
}

// MaybeInt64
type MaybeInt64 struct {
	MaybeError
	value int64
}

func (this *MaybeInt64) Value(value int64) {
	this.Error(nil)
	this.value = value
}

func (this MaybeInt64) Right() int64 {
	this.Test()
	return this.value
}

// MaybeFloat32
type MaybeFloat32 struct {
	MaybeError
	value float32
}

func (this *MaybeFloat32) Value(value float32) {
	this.Error(nil)
	this.value = value
}

func (this *MaybeFloat32) Right() float32 {
	this.Test()
	return this.value
}

// MaybeFloat64
type MaybeFloat64 struct {
	MaybeError
	value float64
}

func (this *MaybeFloat64) Value(value float64) {
	this.Error(nil)
	this.value = value
}

func (this MaybeFloat64) Right() float64 {
	this.Test()
	return this.value
}

// MaybeString
type MaybeString struct {
	MaybeError
	value string
}

func (this *MaybeString) Value(value string) {
	this.Error(nil)
	this.value = value
}

func (this MaybeString) Right() string {
	this.Test()
	return this.value
}

// MaybeBytes
type MaybeBytes struct {
	MaybeError
	value []byte
}

func (this *MaybeBytes) Value(value []byte) {
	this.Error(nil)
	this.value = value
}

func (this MaybeBytes) Right() []byte {
	this.Test()
	return this.value
}

// MaybeEface
type MaybeEface struct {
	MaybeError
	value interface{}
}

func (this *MaybeEface) Value(value interface{}) {
	this.Error(nil)
	this.value = value
}

func (this MaybeEface) Right() interface{} {
	this.Test()
	return this.value
}
