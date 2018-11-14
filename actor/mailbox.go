package actor

import (
	"../message"
	"../config"
	"fmt"
	"errors"
	"../common/maybe"
)

type mailBox struct{
	mailbox chan message.Message
}

func (this *mailBox) Init(attrs interface{}, cfg config.Config) (err maybe.MaybeError) {
	if attrs == nil{
		err.Error(errors.New("actor attrs is nil"))
		return
	}
	attrsMap, ok := attrs.(map[string]interface{})
	if !ok{
		err.Error(fmt.Errorf("illegal cfg type when new actor %s", defaultActorClassName))
		return
	}
	size, ok := attrsMap["MailBoxSize"]
	if !ok{
		err.Error(fmt.Errorf("no actor attribute found: %s", "MailBoxSize"))
		return
	}
	sizeInt, ok := size.(int)
	if !ok {
		err.Error(fmt.Errorf("actor mailbox size cfg type error(expecting int): %+v", size))
		return
	}
	if sizeInt <= 0 {
		err.Error(fmt.Errorf("illegal actor mailbox size: %d", sizeInt))
		return
	}
	this.mailbox = make(chan message.Message, sizeInt)

	return
}
