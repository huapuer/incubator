package protocal

import (
	"incubator/serialization"
	"incubator/common/maybe"
)

type Protocal interface {
	Marshal(serialization.Serializable) []byte
	Unmarshal([]byte, serialization.Serializable) maybe.MaybeError
	GetPackageLen([]byte) maybe.MaybeInt
}
