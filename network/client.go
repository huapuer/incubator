package network

import (
	"../common/maybe"
	"../message"
)

type Client interface {
	Start(string) maybe.MaybeError
	Send(message message.Message) maybe.MaybeError
}
