package serialization

import (
	"../common/maybe"
	"../message"
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

//go:noescape
func Marshal(obj Serializable) (ret []byte) {
	mi := (*mimicIFace)(unsafe.Pointer(&obj))

	size := obj.GetSize()
	ms := &mimicSlice{mi.data, int(size), int(size)}
	val := *(*[]byte)(unsafe.Pointer(ms))

	jbytes := obj.GetJsonBytes().Right()

	ret = append(ret, val...)
	if len(jbytes) > 0 {
		ret = append(ret, jbytes...)
	}

	return
}

//go:noescape
func Unmarshal(data []byte, obj Serializable) (err maybe.MaybeError) {
	lth := int32(len(data))
	lval := obj.GetSize()

	val := data[:lval]
	ms := (*mimicSlice)(unsafe.Pointer(&val))
	mi := (*mimicIFace)(unsafe.Pointer(&obj))
	mi.data = ms.addr

	ljsn := lth - lval
	if ljsn > 0 {
		jsn := data[lval+2 : lth]
		obj.SetJsonField(jsn).Test()
	}

	return
}

//go:noescape
func UnmarshalRemoteMessage(data []byte, msg message.RemoteMessage) (err maybe.MaybeError) {
	lth := int32(len(data))
	lval := msg.GetSize()

	val := data[:lval]
	ms := (*mimicSlice)(unsafe.Pointer(&val))
	mi := (*mimicIFace)(unsafe.Pointer(&msg))
	mi.data = ms.addr

	ljsn := lth - lval
	if ljsn > 0 {
		jsn := data[lval+2 : lth]
		msg.SetJsonField(jsn).Test()
	}

	return
}
