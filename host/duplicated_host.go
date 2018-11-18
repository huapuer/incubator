package host

import (
	"github.com/incubator/common/maybe"
	"../message"
	"errors"
)

type duplicatedHost struct{
	commonHost

	master Host
	slaves []Host
}

func NewDuplicatedHost(master Host, slaves []Host) (ret MaybeHost) {
	ret.Value(&duplicatedHost{
		master:master,
		slaves:slaves,
	})
	return
}

func (this *duplicatedHost) Receive(msg message.Message) (err maybe.MaybeError) {
	if this.master == nil {
		err.Error(errors.New("master host not set"))
		return
	}
	msg.SetHostId(this.GetId().Right()).Test()
	this.master.Receive(msg).Test()
	for _, slave := range this.slaves {
		if slave == nil {
			err.Error(errors.New("nil slave host found"))
			return
		}
		slave.Receive(msg).Test()
	}
	return
}

func (this duplicatedHost) GetJsonBytes() (ret maybe.MaybeBytes) {
	panic("not implemented")
}

func (this *duplicatedHost) SetJsonField(data []byte) (err maybe.MaybeError) {
	panic("not implemented")
}

func (this duplicatedHost) GetSize() int32 {
	panic("not implemented")
}