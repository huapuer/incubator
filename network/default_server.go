package network

import (
	"../common/maybe"
	"../message"
	"errors"
	"incubator/protocal"
)

type defaultServer struct {
	commonServer

	p protocal.Protocal
}

//go:noescape
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

//go:noescape
func (this defaultServer) handleData(data []byte, l int) (err maybe.MaybeError) {
	if l == 0 {
		err.Error(errors.New("empty data"))
		return
	}
	this.packageBuffer = append(this.packageBuffer, data...)
	if this.packageSize < 0 {
		this.packageSize = this.p.GetPackageLen(data).Right()
	}
	if this.packageSize >= 0{
		if this.readBufferSize >= this.packageSize {
			pkg := this.packageBuffer[:this.packageSize]
			this.packageBuffer = this.packageBuffer[this.packageSize:]
			this.handlePackage(pkg).Test()
		}
	}

	return
}
