package network

import (
	"../common/maybe"
	"../message"
	"../config"
	"fmt"
	"errors"
	"incubator/serialization"
	"time"
)

const (
	DefaultClient = &defaultClient{}
)

type defaultClient struct {
	config.IOC

	pool connectionPool
}

type MaybeDefualtClient struct {
	config.IOC

	maybe.MaybeError
	value *defaultClient
}

func (this MaybeDefualtClient) Value(value *defaultClient) {
	this.Error(nil)
	this.value = value
}

func (this MaybeDefualtClient) Right() *defaultClient {
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
		return ret
	}
	attrsMap, ok := attrs.(map[string]interface{})
	if !ok {
		ret.Error(errors.New("illegal attrs type when new default client"))
		return ret
	}

	maxIdle, ok := attrsMap["MaxSpare"]
	if !ok {
		ret.Error(errors.New("attribute MaxSpare not found"))
		return ret
	}
	maxIdleInt, ok := maxIdle.(int32)
	if !ok {
		ret.Error(fmt.Errorf("max idle cfg type error(expecting int): %+v", maxIdle))
		return ret
	}

	maxBusy, ok := attrsMap["MaxBusy"]
	if !ok {
		ret.Error(errors.New("attribute MaxBusy not found"))
		return ret
	}
	maxBusyInt, ok := maxBusy.(int32)
	if !ok {
		ret.Error(fmt.Errorf("max busy cfg type error(expecting int): %+v", maxIdle))
		return ret
	}

	timeout, ok := attrsMap["Timeout"]
	if !ok {
		ret.Error(errors.New("attribute Timeout not found"))
		return ret
	}
	timeoutInt, ok := timeout.(int64)
	if !ok {
		ret.Error(fmt.Errorf("timeout cfg type error(expecting int): %+v", timeout))
		return ret
	}

	addr, ok := attrsMap["Address"]
	if !ok {
		ret.Error(errors.New("attribute Address not found"))
		return ret
	}
	addrStr, ok := timeout.(string)
	if !ok {
		ret.Error(fmt.Errorf("address cfg type error(expecting string): %+v", addr))
		return ret
	}

	ret.Value(&defaultClient{
		pool: NewConnectionPool(addrStr, maxIdleInt, maxBusyInt, time.Duration(timeoutInt)).Right(),
	})
	return ret
}

func (this *defaultClient) Send(msg message.Message) (err maybe.MaybeError) {
	_, e := this.pool.GetConnection().Right().Write(serialization.Marshal(msg))
	if e != nil {
		err.Error(e)
	}
	return
}
