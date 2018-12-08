package network

import (
	"../common/maybe"
	"../config"
	"../message"
	"github.com/incubator/protocal"
	"time"
)

var (
	DefaultClient = &defaultClient{}
)

type defaultClient struct {
	config.IOC

	maxIdle int32
	maxBusy int32
	timeout int64
	pool    connectionPool
	p       protocal.Protocal
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

	this.maxIdle = config.GetAttrInt32(attrs, "MaxIdle", config.CheckInt32GT0).Right()
	this.maxBusy = config.GetAttrInt32(attrs, "MaxBusy", config.CheckInt32GT0).Right()
	this.timeout = config.GetAttrInt64(attrs, "Timeout", config.CheckInt64GT0).Right()

	protocalStr := config.GetAttrString(attrs, "Protocal", config.CheckStringNotEmpty).Right()
	protocal := protocal.GetProtocalPrototype(protocalStr).Right()

	ret.Value(&defaultClient{
		maxIdle: config.GetAttrInt32(attrs, "MaxIdle", config.CheckInt32GT0).Right(),
		maxBusy: config.GetAttrInt32(attrs, "MaxBusy", config.CheckInt32GT0).Right(),
		timeout: config.GetAttrInt64(attrs, "Timeout", config.CheckInt64GT0).Right(),
		p:       protocal,
	})
	return ret
}

func (this *defaultClient) Connect(addr string) {
	this.pool = NewConnectionPool(addr, this.maxIdle, this.maxBusy, time.Duration(this.timeout)).Right()
}

//go:noescape
func (this *defaultClient) Send(msg message.RemoteMessage) (err maybe.MaybeError) {
	_, e := this.pool.GetConnection().Right().Write(this.p.Pack(msg))
	if e != nil {
		err.Error(e)
	}
	return
}
