package protocal

import (
	"github.com/incubator/interfaces"
	"github.com/incubator/serialization"
	"unsafe"
)

const (
	fixedHeaderProtocalClassName = "protocal.fixedHeaderProtocal"
)

func init() {
	RegisterProtocalPrototype(fixedHeaderProtocalClassName, &fixedHeaderProtocal{}).Test()
}

type fixedHeaderProtocal struct{}

func (this fixedHeaderProtocal) New(attrs interface{}, cfg interfaces.Config) interfaces.IOC {
	ret := interfaces.MaybeProtocal{}
	ret.Value(&fixedHeaderProtocal{})
	return ret
}

////go:noescape
func (this *fixedHeaderProtocal) Pack(msg interfaces.RemoteMessage) (ret []byte) {
	bytes := serialization.Marshal(msg)
	lth := len(bytes) + int(unsafe.Sizeof(int32(0)))

	ret = serialization.Ptr2Bytes(unsafe.Pointer(&lth), int(unsafe.Sizeof(int32(0))))
	ret = append(ret, bytes...)

	return
}

func (this *fixedHeaderProtocal) Parse(data []byte) (int, int) {
	l := len(data)
	if l < 4 {
		return PROTOCAL_PARSE_STATE_SHORT, 0
	}
	var lth int
	serialization.MoveBytes(unsafe.Pointer(&lth), data, int(unsafe.Sizeof(int32(0))))
	if lth < 0 {
		return PROTOCAL_PARSE_STATE_ERROR, 0
	}
	if l < lth {
		return PROTOCAL_PARSE_STATE_SHORT, 0
	}

	return lth, int(unsafe.Sizeof(int32(0)))
}

func (this *fixedHeaderProtocal) Decode(data []byte) (ret []byte) {
	ret = data
	return
}
