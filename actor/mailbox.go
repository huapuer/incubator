package actor

import (
	"../common/maybe"
	"../config"
	"../message"
)

type mailBox struct {
	mailbox chan message.Message
}

func (this *mailBox) Init(attrs interface{}, cfg config.Config) (err maybe.MaybeError) {
	size := config.GetAttrInt(attrs, "MailBoxSize", config.CheckIntGT0).Right()
	this.mailbox = make(chan message.Message, size)

	return
}
