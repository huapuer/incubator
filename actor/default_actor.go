package actor

import (
	c "../context"
	"context"
	"errors"
	"fmt"
	"../common/maybe"
	"../config"
	"../message"
)

const (
	defaultActorClassName = "actor.defaultActor"
)

func init() {
	RegisterActorPrototype(defaultActorClassName, &defaultActor{}).Test()
}

type defaultActor struct {
	commonActor

	mailbox chan message.Message
}

func (this defaultActor) New(attrs interface{}, cfg config.Config) config.IOC {
	ret := MaybeActor{}
	if attrs == nil{
		ret.Error(errors.New("actor attrs is nil"))
		return ret
	}
	attrsMap, ok := attrs.(map[string]interface{})
	if !ok{
		ret.Error(fmt.Errorf("illegal cfg type when new actor %s", defaultActorClassName))
		return ret
	}
	size, ok := attrsMap["MailBoxSize"]
	if !ok{
		ret.Error(fmt.Errorf("no actor attribute found: %s", "MailBoxSize"))
		return ret
	}
	sizeInt, ok := size.(int)
	if !ok {
		ret.Error(fmt.Errorf("actor mailbox size cfg type error(expecting int): %+v", size))
		return ret
	}
	if sizeInt <= 0 {
		ret.Error(fmt.Errorf("illegal actor mailbox size: %d", sizeInt))
		return ret
	}
	ret.Value(&defaultActor{mailbox: make(chan message.Message, sizeInt)})

	return ret
}

func (this *defaultActor) Start(ctx context.Context) (err maybe.MaybeError) {
	if this.mailbox == nil {
		err.Error(errors.New("mailbox not inited."))
		return
	}

	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				err.Error(errors.New("ctx done."))
				return
			case m := <-this.mailbox:
				maybe.TryCatch(
					func() {
						ctx := c.MessageContext{
							Topo: this.Topo,
						}
						m.Process(ctx).Test()
					}, nil)
			}
		}
	}(ctx)

	return
}

func (this *defaultActor) Receive(msg message.Message) (err maybe.MaybeError) {
	if this.mailbox == nil {
		err.Error(errors.New("mailbox not inited."))
		return
	}

	this.mailbox <- msg
	return
}
