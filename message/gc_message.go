package message

import (
	"github.com/incubator/actor"
	"github.com/incubator/common/maybe"
	"runtime"
	"time"
	"unsafe"
)

const (
	GCMessageClassName = "message.GCMessage"
)

func init() {
	RegisterMessagePrototype(GCMessageClassName, &GCMessage{
		commonMessage: commonMessage{
			layer:  -1,
			typ:    -1,
			master: -1,
			hostId: -1,
		},
	}).Test()
}

type GCMessage struct {
	commonMessage

	interval time.Duration
}

func (this *GCMessage) Process(runner actor.Actor) (err maybe.MaybeError) {
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
		func(runner actor.Actor) {
			runner.Receive(this)
		})

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

func (this *GCMessage) Replicate() (ret MaybeRemoteMessage) {
	new := *this
	ret.Value(&new)
	return
}
