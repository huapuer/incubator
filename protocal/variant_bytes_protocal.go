package protocal

import (
	"../config"
	"../message"
	"../serialization"
)

const (
	variantBytesProtocalClassName = "protocal.variantBytesProtocal"
)

func init() {
	RegisterProtocalPrototype(variantBytesProtocalClassName, &variantBytesProtocal{}).Test()
}

type variantBytesProtocal struct{}

func (this variantBytesProtocal) New(attrs interface{}, cfg config.Config) config.IOC {
	ret := MaybeProtocal{}
	ret.Value(&fixedHeaderProtocal{})
	return ret
}

//go:noescape
func (this *variantBytesProtocal) Pack(msg message.RemoteMessage) (ret []byte) {
	bytes := serialization.Marshal(msg)

	var (
		l   = len(bytes)
		s   = uint(1)
		idx = 0
	)

	for i := 0; i < l+1; {
		b := uint8(0)
		if i > 0 {
			sp := 8 - s
			b |= (bytes[i-1] & (0xFF >> sp)) << sp
		}
		if i < l {
			b = (bytes[i] & (0xFF << s)) >> s
		}

		ret = append(ret, b)
		idx++
		if idx > 0 {
			ret[idx-1] |= 0x80
		}

		s++
		if s == 8 {
			s = 1
		} else {
			i++
		}
	}

	return
}

func (this *variantBytesProtocal) Parse(data []byte) (int, int) {
	for i := len(data) - 1; i >= 0; i-- {
		if data[i]&0x80 == 0 {
			return i + 1, 0
		}
	}

	return PROTOCAL_PARSE_STATE_SHORT, 0
}

func (this *variantBytesProtocal) Decode(data []byte) (ret []byte) {
	var (
		l = len(data)
		s = uint(1)
	)

	for i := 0; i < l; {
		b := (data[i] & (0xFF >> s)) << s
		if i < l-1 {
			sp := 8 - s
			b = (data[i+1] & (0xFF << sp)) >> sp
		}

		ret = append(ret, b)

		s++
		if s == 8 {
			s = 1
		} else {
			i++
		}
	}

	return
}
