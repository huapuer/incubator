package protocal

import (
	"unsafe"
	"incubator/common/maybe"
	"fmt"
	"../serialization"
)

type fixHeaderProtocal struct {}

//go:noescape
func (this *fixHeaderProtocal) Marshal(obj serialization.Serializable) (ret []byte) {
	jbytes := obj.GetJsonBytes().Right()
	lth := obj.GetSize() + int32(len(jbytes) + 1 + int(unsafe.Sizeof(int32(0))))

	ret = append(ret, uint8(lth))
	ret = append(ret, serialization.Marshal(obj)...)

	return
}

//go:noescape
func (this *fixHeaderProtocal) Unmarshal(data []byte, obj serialization.Serializable) (err maybe.MaybeError) {
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

	val := data[2:]

	return  serialization.Unmarshal(val, obj)
}

func (this *fixHeaderProtocal) GetPackageLen(data []byte) (ret maybe.MaybeInt) {
	l := len(data)
	if l < 4 {
		ret.Value(-1)
		return
	}
	lth := int(data[0])
	if lth < 0 {
		ret.Error(fmt.Errorf("message claims negative lenth: %d", lth))
		return
	}
	if l != lth {
		ret.Error(fmt.Errorf("message length not equal to claimed, %d != %d", l, lth))
		return
	}

	ret.Value(lth)
	return
}
