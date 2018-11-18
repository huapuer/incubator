package serialization

import (
	"github.com/incubator/common/maybe"
	"fmt"
	"unsafe"
)

type Serializable interface {
	GetSize() int32
	GetJsonBytes() maybe.MaybeBytes
	SetJsonField([]byte) maybe.MaybeError
}

type MaybeSerializable struct {
	maybe.MaybeError
	value Serializable
}

func (this MaybeSerializable) Value(value Serializable) {
	this.Error(nil)
	this.value = value
}

func (this MaybeSerializable) Right() Serializable {
	this.Test()
	return this.value
}

type mimicSlice struct {
	addr *unsafe.ArbitraryType
	len  int
	cap  int
}

type mimicIFace struct {
	tab  unsafe.Pointer
	data unsafe.Pointer
}

func Marshal(obj Serializable) (ret []byte) {
	mi := (*mimicIFace)(unsafe.Pointer(&obj))

	size := obj.GetSize()
	ms := &mimicSlice{mi.data, size, size}
	val := *(*[]byte)(unsafe.Pointer(ms))

	jbytes := obj.GetJsonBytes().Right()

	lth := int32(len(val) + len(jbytes) + 1 + unsafe.Sizeof(int32(0)))

	ret = append(ret, uint8(lth), uint8(len(val)))
	ret = append(ret, val...)
	if len(jbytes) > 0 {
		ret = append(ret, jbytes...)
	}

	return
}

func Unmarshal(data []byte, obj Serializable) (err maybe.MaybeError) {
	l := len(data)
	if l < 4 {
		err.Error(fmt.Errorf("message bytes too short: %d", l))
		return
	}
	lth := int(data[0])
	if lth < 0 {
		err.Error(fmt.Errorf("message claims negative lenth: %d", lth))
		return
	}
	if l != lth {
		err.Error(fmt.Errorf("message length not equal to claimed, %d != %d", l, lth))
		return
	}

	lval := int(data[1])
	if lth < 0 {
		err.Error(fmt.Errorf("message claims negative binary length: %d", lval))
		return
	}
	if lth < lval+3 {
		err.Error(fmt.Errorf("message length shorter than claimed binary length + header length, %d  < %d + 3", l, lval))
		return
	}

	val := data[2 : lval+2]

	ms := (*mimicSlice)(unsafe.Pointer(&val))

	mi := (*mimicIFace)(unsafe.Pointer(&obj))
	mi.data = ms.addr

	ljsn := lth - lval - 2
	if ljsn > 0 {
		jsn := data[lval+2 : lth]
		obj.SetJsonField(jsn).Test()
	}

	return
}