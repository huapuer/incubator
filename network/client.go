package network

import (
	"../common/maybe"
	"../message"
	"github.com/incubator/config"
)

type Client interface {
	config.IOC

	Send(message message.Message) maybe.MaybeError
}
