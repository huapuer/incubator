package network

import (
	"../common/maybe"
	"../config"
	"../message"
)

type Client interface {
	config.IOC

	Send(message message.RemoteMessage) maybe.MaybeError
}
