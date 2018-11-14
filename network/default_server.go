package network

import (
	"errors"
	"../common/maybe"
	"../message"
)

type defaultServer struct {
	commonServer
}

func (this defaultServer) handlePacakge(data []byte) (err maybe.MaybeError) {
	if len(data) <= 2 {
		err.Error(errors.New("empty package"))
		return
	}
	layer := data[0]
	typ := data[1]
	message.RoutePackage(data, uint8(layer), uint8(typ)).Test()
	return
}
