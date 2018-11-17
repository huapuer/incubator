package network

import (
	"../common/maybe"
	"../message"
	"github.com/incubator/config"
	"fmt"
	"errors"
)

const (
	DefaultClient = defaultClient{}
)

type defaultClient struct {
	config.IOC

	pool connectionPool
}

type MaybeDefualtClient struct {
	config.IOC

	maybe.MaybeError
	value defaultClient
}

func (this MaybeDefualtClient) Value(value defaultClient) {
	this.Error(nil)
	this.value = value
}

func (this MaybeDefualtClient) Right() defaultClient {
	this.Test()
	return this.value
}

func (this MaybeDefualtClient) New(cfg config.Config, args ...int32) config.IOC {
	panic("not implemented.")
}

func (this *defaultClient) New(attrs interface{}, cfg config.Config) config.IOC {
	ret := MaybeDefualtClient{}

	if attrs == nil {
		ret.Error(errors.New("attrs is nil when new default client"))
		return
	}
	attrsMap, ok := attrs.(map[string]interface{})
	if !ok {
		ret.Error(errors.New("illegal attrs type when new default client"))
		return
	}

	maxIdle, ok := attrsMap["MaxSpare"]
	if !ok {
		ret.Error(errors.New("attribute MaxSpare not found"))
		return
	}
	maxIdleInt, ok := maxIdle.(int32)
	if !ok {
		ret.Error(fmt.Errorf("max idle cfg type error(expecting int): %+v", maxIdle))
		return
	}

	maxBusy, ok := attrsMap["MaxBusy"]
	if !ok {
		ret.Error(errors.New("attribute MaxBusy not found"))
		return
	}
	maxBusyInt, ok := maxBusy.(int32)
	if !ok {
		ret.Error(fmt.Errorf("max busy cfg type error(expecting int): %+v", maxIdle))
		return
	}

	timeout, ok := attrsMap["Timeout"]
	if !ok {
		ret.Error(errors.New("attribute Timeout not found"))
		return
	}
	timeoutInt, ok := timeout.(int64)
	if !ok {
		ret.Error(fmt.Errorf("timeout cfg type error(expecting int): %+v", timeout))
		return
	}

	addr, ok := attrsMap["Address"]
	if !ok {
		ret.Error(errors.New("attribute Address not found"))
		return
	}
	addrStr, ok := timeout.(string)
	if !ok {
		ret.Error(fmt.Errorf("address cfg type error(expecting string): %+v", addr))
		return
	}

	ret.Value(&defaultClient{
		pool: NewConnectionPool(addrStr, maxIdleInt, maxBusyInt, timeoutInt).Right(),
	})
	return
}

func (this *defaultClient) Send(msg message.Message) (err maybe.MaybeError) {
	_, e := this.pool.GetConnection().Right().Write(msg.Marshal(msg))
	if err != nil {
		err.Error(e)
	}
	return
}
