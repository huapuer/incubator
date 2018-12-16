package network

import (
	"github.com/incubator/common/maybe"
	"github.com/incubator/config"
	"github.com/incubator/interfaces"
	"github.com/incubator/protocal"
	"time"
)

const (
	defaultClientClassName = "network.defaultClient"
)

func init() {
	RegisterClientPrototype(defaultClientClassName, &defaultClient{}).Test()
}

var (
	DefaultClient = &defaultClient{}
)

type defaultClient struct {
	interfaces.IOC

	maxIdle int32
	maxBusy int32
	timeout int64
	pool    connectionPool
	p       interfaces.Protocal
}

func (this *defaultClient) New(attrs interface{}, cfg interfaces.Config) interfaces.IOC {
	ret := interfaces.MaybeClient{}

	maxIdle := config.GetAttrInt32(attrs, "MaxIdle", config.CheckInt32GT0).Right()
	maxBusy := config.GetAttrInt32(attrs, "MaxBusy", config.CheckInt32GT0).Right()
	timeout := config.GetAttrInt64(attrs, "Timeout", config.CheckInt64GT0).Right()

	protocalStr := config.GetAttrString(attrs, "Protocal", config.CheckStringNotEmpty).Right()
	protocal := protocal.GetProtocalPrototype(protocalStr).Right()

	ret.Value(&defaultClient{
		maxIdle: maxIdle,
		maxBusy: maxBusy,
		timeout: timeout,
		p:       protocal,
	})
	return ret
}

func (this *defaultClient) Connect(addr string) {
	this.pool = NewConnectionPool(addr, this.maxIdle, this.maxBusy, time.Duration(this.timeout)).Right()
}

////go:noescape
func (this *defaultClient) Send(msg interfaces.RemoteMessage) (err maybe.MaybeError) {
	_, e := this.pool.GetConnection().Right().Write(this.p.Pack(msg))
	if e != nil {
		err.Error(e)
	}
	return
}
