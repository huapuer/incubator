// +build debug

package maybe

import (
	"fmt"
	"runtime"
	//"github.com/sirupsen/logrus"
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
			if catch != nil {
				catch(err.(error))
			}
		}
	}()
	try()
}

func getFormattedCaller() string {
	fpcs := make([]uintptr, 1)

	n := runtime.Callers(3, fpcs)
	if n == 0 {
		return "n/a"
	}

	fun := runtime.FuncForPC(fpcs[0] - 1)
	if fun == nil {
		return "n/a"
	}

	file, line := fun.FileLine(fpcs[0] - 1)

	return fmt.Sprintf("%s(%d): %s", file, line, fun.Name())
}

type Maybe interface {
	Error(error)
	Test()
}

type MaybeError struct {
	nonNil bool
	err    error
}

type StackInfoErr struct {
	stackInfo string
	e         error
}

func (this StackInfoErr) Error() string {
	return fmt.Sprintf("%s:%s", this.stackInfo, this.e)
}

func (this *MaybeError) Error(err error) {
	if err == nil {
		this.nonNil = true
		return
	}

	this.err = &StackInfoErr{
		stackInfo: getFormattedCaller(),
		e:         err,
	}
}

func (this MaybeError) Test() {
	if this.err != nil {
		log.Printf("[debugger]%s", this.err)
		panic(this.err)
	}

	if this.nonNil == false {
		err := NilValueError{"Value not set."}
		log.Printf("[debugger]%s, %+v", getFormattedCaller(), err)
		panic(err)
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
	log.Printf("[debugger]%s, %d", getFormattedCaller(), this.value)
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
	log.Printf("[debugger]%s, %f", getFormattedCaller(), this.value)
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
	log.Printf("[debugger]%s, %f", getFormattedCaller(), this.value)
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
	log.Printf("[debugger]%s, %s", getFormattedCaller(), this.value)
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
	log.Printf("[debugger]%s, %s", getFormattedCaller(), string(this.value))
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
	log.Printf("[debugger]%s, %+v", getFormattedCaller(), this.value)
	return this.value
}
