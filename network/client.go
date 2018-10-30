package network

import (
	"incubator/common/maybe"
	"incubator/message"
)

type Client interface {
	Start(string) maybe.MaybeError
	Send(message message.Message) maybe.MaybeError
}
