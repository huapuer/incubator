package io

import (
	"github.com/incubator/interfaces"
)

type joint struct {
	begin  int64
	end    int64
	client interfaces.Client
}

type commonIO struct {
	layerId int32
}

func (this *commonIO) SetLayer(id int32) {
	this.layerId = id
}
