package network

import (
	"errors"
	"../common/maybe"
	"../message"
)

type defaultServer struct {
	commonServer
}

func (this defaultServer) handlePackage(data []byte) (err maybe.MaybeError) {
	if len(data) <= 2 {
		err.Error(errors.New("empty package"))
		return
	}
	layer := data[0]
	typ := data[1]
	message.RoutePackage(data, uint8(layer), uint8(typ)).Test()
	return
}


func (this defaultServer) handleData(data []byte, l int) (err maybe.MaybeError) {
	if l == 0 {
		err.Error(errors.New("empty data"))
		return
	}
	if this.packageSize == 0 {
		this.packageSize = int(data[0])
		this.packageBuffer = data[1:]
	} else {
		want := len(this.packageBuffer) + l - this.packageSize
		if want >= 0 {
			pkg := this.packageBuffer
			pkg = append(pkg, data[:want]...)
			if want > 0 {
				this.packageBuffer = data[want:]
			}
			this.packageSize = 0
			this.handlePackage(pkg).Test()
		} else {
			this.packageBuffer = append(this.packageBuffer, data...)
		}
	}

	return
}
