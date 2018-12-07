package router

import (
	"../actor"
	"../common/maybe"
	"../config"
	"../message"
	"runtime"
)

const (
	spikeRouterClassName = "router.defaultRouter"
)

func init() {
	RegisterRouterPrototype(spikeRouterClassName, &defaultRouter{}).Test()
}

type spikeRouter struct {
	router defaultRouter
}

func (this spikeRouter) New(attrs interface{}, cfg config.Config) config.IOC {
	ret := MaybeRouter{}

	maybe.TryCatch(
		func() {
			r := defaultRouter{}.New(attrs, cfg).(MaybeRouter).Right()
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

//go:noescape
func (this spikeRouter) Route(msg message.RemoteMessage) (err maybe.MaybeError) {
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

func (this spikeRouter) GetActors() []actor.Actor {
	return this.router.GetActors()
}

func (this spikeRouter) Stop() {
	this.router.Stop()
}
