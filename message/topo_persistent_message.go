package message

import (
	"fmt"
	"github.com/incubator/common/maybe"
	"github.com/incubator/interfaces"
	"github.com/incubator/persistence"
	"time"
	"unsafe"
)

const (
	TopoPersistentMessageClassName = "message.TopoPersistentMessage"
)

func init() {
	interfaces.RegisterMessagePrototype(TopoPersistentMessageClassName, &TopoPersistentMessage{
		commonMessage: commonMessage{
			layerId: -1,
			typ:     -1,
			master:  -1,
			hostId:  -1,
		},
	}).Test()
}

type TopoPersistentMessage struct {
	commonMessage

	layer           int32
	storeExpiration time.Duration
	loadExpiration  time.Duration
	interval        time.Duration
}

func (this *TopoPersistentMessage) Process(runner interfaces.Actor) (err maybe.MaybeError) {
	runner.SetState(runner, fmt.Sprintf("topo_%d_persistent_touch", this.layer), true, this.interval,
		func(runner interfaces.Actor) {
			runner.Receive(this)
		}).Test()

	pa, ok := interfaces.GetLayer(this.layer).Right().GetTopo().(persistence.Persistentable)
	if !ok {
		err.Error(fmt.Errorf("topo of layer(%d) is not persistentable", this.layer))
		return
	}

	pa.SetStoreExpiration(this.storeExpiration)
	pa.SetLoadExpiration(this.loadExpiration)

	pa.Persistent().Test()

	return
}

func (this *TopoPersistentMessage) GetJsonBytes() (ret maybe.MaybeBytes) {
	ret.Error(nil)
	return
}

func (this *TopoPersistentMessage) SetJsonField(data []byte) (err maybe.MaybeError) {
	err.Error(nil)
	return
}

func (this *TopoPersistentMessage) GetSize() int32 {
	return int32(unsafe.Sizeof(*this))
}
