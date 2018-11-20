package network

import (
	"../common/maybe"
	"../message"
	"../config"
)

type Client interface {
	config.IOC

	Send(message message.Message) maybe.MaybeError
}
