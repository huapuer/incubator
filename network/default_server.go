package network

import (
	"errors"
	"incubator/common/maybe"
	"incubator/message"
)

type defaultServer struct {
	commonServer
}

func (this defaultServer) handlePacakge(data []byte) (err maybe.MaybeError) {
	if len(data) <= 1 {
		err.Error(errors.New("empty package"))
		return
	}
	typ := data[0]
	message.RoutePackage(data, int(typ)).Test()
	return
}
