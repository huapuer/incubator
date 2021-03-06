package message

import (
	"github.com/incubator/common/maybe"
	"github.com/incubator/interfaces"
	"runtime"
	"time"
	"unsafe"
)

const (
	GCMessageClassName = "message.GCMessage"
)

func init() {
	interfaces.RegisterMessagePrototype(GCMessageClassName, &GCMessage{
		commonMessage: commonMessage{
			layerId: -1,
			typ:     -1,
			master:  -1,
			hostId:  -1,
		},
	}).Test()
}

type GCMessage struct {
	commonMessage

	interval time.Duration
}

func (this *GCMessage) Process(runner interfaces.Actor) (err maybe.MaybeError) {
	started := false
	maybe.TryCatch(
		func() {
			started = runner.GetState("gc_started").Right().(bool)
		}, nil)
	if started {
		err.Error(nil)
		return
	}
	runner.SetState(runner, "gc_started", true, this.interval,
		func(runner interfaces.Actor) {
			runner.Receive(this)
		}).Test()

	//ms := runtime.MemStats{}
	//runtime.ReadMemStats(&ms)

	runtime.GC()

	return
}

func (this *GCMessage) GetJsonBytes() (ret maybe.MaybeBytes) {
	ret.Error(nil)
	return
}

func (this *GCMessage) SetJsonField(data []byte) (err maybe.MaybeError) {
	err.Error(nil)
	return
}

func (this *GCMessage) GetSize() int32 {
	return int32(unsafe.Sizeof(*this))
}
