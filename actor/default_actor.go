package actor

import (
	"context"
	"errors"
	"fmt"
	"incubator/common/maybe"
	"incubator/config"
	"incubator/message"
	"strconv"
)

const (
	defaultActorClassName = "actor.defaultActor"
)

func init() {
	RegisterActorPrototype(defaultActorClassName, &defaultActor{}).Test()
}

type defaultActor struct {
	mailbox chan message.Message
}

func (this defaultActor) New(cfg config.Config) config.IOC {
	ret := MaybeActor{}
	if actor, ok := cfg.Actors[defaultActorClassName]; ok {
		attrs := actor.Attributes.(map[string]string)
		if size, ok := attrs["MailBoxSize"]; ok {
			if mailBoxSize, err := strconv.Atoi(size); err == nil {
				return NewDefaultActor(int64(mailBoxSize))
			}
			ret.Error(fmt.Errorf("illegal actor attribute: %s=%s", "MailBoxSize", size))
			return ret
		}
		ret.Error(fmt.Errorf("no actor attribute found: %s", "MailBoxSize"))
		return ret
	}
	ret.Error(fmt.Errorf("no actor class cfg found: %s", defaultActorClassName))
	return ret
}

func NewDefaultActor(taskChanSize int64) (this MaybeActor) {
	if taskChanSize <= 0 {
		this.Error(fmt.Errorf("wrong task chan size: %d", taskChanSize))
		return
	}
	this.Value(&defaultActor{make(chan message.Message, taskChanSize)})
	return
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
