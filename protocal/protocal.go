package protocal

import (
	"github.com/incubator/message"
)

const (
	PROTOCAL_PARSE_STATE_SHORT = -1
	PROTOCAL_PARSE_STATE_ERROR = -2
)

type Protocal interface {
	Pack(message.RemoteMessage) []byte
	Parse([]byte) (int, int)
}
