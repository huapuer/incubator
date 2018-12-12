package router

import (
	"github.com/incubator/common/maybe"
	"github.com/incubator/interfaces"
	"runtime"
)

const (
	spikeRouterClassName = "router.spikeRouter"
)

func init() {
	RegisterRouterPrototype(spikeRouterClassName, &defaultRouter{}).Test()
}

type spikeRouter struct {
	router defaultRouter
}

func (this spikeRouter) New(attrs interface{}, cfg interfaces.Config) interfaces.IOC {
	ret := interfaces.MaybeRouter{}

	maybe.TryCatch(
		func() {
			r := defaultRouter{}.New(attrs, cfg).(interfaces.MaybeRouter).Right()
			ret.Value(&spikeRouter{r.(defaultRouter)})
		},
		func(err error) {
			ret.Error(err)
		})

	return ret
}

func (this spikeRouter) Start() {
	this.router.Start()
}

////go:noescape
func (this spikeRouter) Route(msg interfaces.RemoteMessage) (err maybe.MaybeError) {
	maybe.TryCatch(
		func() {
			this.router.Route(msg).Test()
			runtime.Gosched()
		},
		func(e error) {
			err.Error(e)
		})

	return
}

func (this spikeRouter) SimRoute(seed int64, actorsNum int) int64 {
	return this.router.SimRoute(seed, actorsNum)
}

func (this spikeRouter) GetActors() []interfaces.Actor {
	return this.router.GetActors()
}

func (this spikeRouter) Stop() {
	this.router.Stop()
}
