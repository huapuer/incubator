package serialization

import (
	"github.com/incubator/common/maybe"
	"github.com/incubator/interfaces"
	"unsafe"
)

////go:noescape
func Marshal(obj interfaces.Serializable) (ret []byte) {
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

////go:noescape
func Unmarshal(data []byte, obj interfaces.Serializable) (err maybe.MaybeError) {
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

////go:noescape
func UnmarshalRemoteMessage(data []byte, msg interfaces.RemoteMessage) (err maybe.MaybeError) {
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
