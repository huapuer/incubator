package actor

import (
	"github.com/incubator/common/maybe"
	"github.com/incubator/config"
	"github.com/incubator/interfaces"
)

type mailBox struct {
	mailbox chan interfaces.Message
}

func (this *mailBox) Init(attrs interface{}, cfg *config.Config) (err maybe.MaybeError) {
	size := config.GetAttrInt(attrs, "MailBoxSize", config.CheckIntGT0).Right()
	this.mailbox = make(chan interfaces.Message, size)

	return
}
