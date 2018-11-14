package actor

import (
	c "../context"
	"context"
	"errors"
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
	mailBox
}

func (this defaultActor) New(attrs interface{}, cfg config.Config) config.IOC {
	ret := MaybeActor{}

	actor := &defaultActor{}
	actor.mailBox.Init(attrs, cfg).Test()

	ret.Value(actor)
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
