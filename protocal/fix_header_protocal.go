package protocal

import (
	"../config"
	"../message"
	"../serialization"
	"unsafe"
)

const (
	fixHeaderProtocalClassName = "protocal.fixHeaderProtocal"
)

func init() {
	RegisterProtocalPrototype(fixHeaderProtocalClassName, &fixHeaderProtocal{}).Test()
}

type fixHeaderProtocal struct{}

func (this fixHeaderProtocal) New(attrs interface{}, cfg config.Config) config.IOC {
	ret := MaybeProtocal{}
	ret.Value(&fixHeaderProtocal{})
	return ret
}

func (this *fixHeaderProtocal) Pack(msg message.RemoteMessage) (ret []byte) {
	bytes := serialization.Marshal(msg)
	lth := len(bytes) + 1 + int(unsafe.Sizeof(int32(0)))

	ret = append(ret, uint8(lth))
	ret = append(ret, bytes...)

	return
}

func (this *fixHeaderProtocal) Parse(data []byte) (int, int) {
	l := len(data)
	if l < 4 {
		return PROTOCAL_PARSE_STATE_SHORT, 0
	}
	lth := int(data[0])
	if lth < 0 {
		return PROTOCAL_PARSE_STATE_ERROR, 0
	}
	if l != lth {
		return PROTOCAL_PARSE_STATE_ERROR, 0
	}

	return lth, int(unsafe.Sizeof(int32(0)))
}
