package network

import (
	"../common/maybe"
	"../config"
	"../message"
)

type Client interface {
	config.IOC

	Connect(string)
	Send(message message.RemoteMessage) maybe.MaybeError
}
