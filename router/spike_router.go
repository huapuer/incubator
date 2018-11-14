package router

import (
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
		func(){
			r := defaultRouter{}.New(attrs, cfg).(MaybeRouter).Right()
			ret.Value(&spikeRouter{r})
		},
		func(err error){
			ret.Error(err)
		})

	return ret
}

func (this spikeRouter) Route(msg message.Message) (err maybe.MaybeError) {
	maybe.TryCatch(
		func(){
			this.router.Route(msg).Test()
			runtime.Gosched()
		},
		func(e error){
			err.Error(e)
		})

	return
}
