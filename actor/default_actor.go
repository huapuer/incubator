package actor

import (
	"context"
	"errors"
	"github.com/incubator/common/maybe"
	"github.com/incubator/config"
	"github.com/incubator/interfaces"
	"time"
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
	defaultHealthManager
}

func (this defaultActor) New(attrs interface{}, cfg interfaces.Config) interfaces.IOC {
	ret := interfaces.MaybeActor{}

	heartbeatIntvl := config.GetAttrInt64(attrs, "HeartbeatIntvl", config.CheckInt64GT0).Right()

	actor := &defaultActor{
		defaultHealthManager: defaultHealthManager{
			heartbeatIntvl: time.Duration(heartbeatIntvl),
		},
	}
	actor.mailBox.Init(attrs, cfg.(*config.Config)).Test()

	ret.Value(actor)
	return ret
}

func (this *defaultActor) Start(ctx context.Context) (err maybe.MaybeError) {
	if this.mailbox == nil {
		err.Error(errors.New("mailbox not inited."))
		return
	}

	this.defaultHealthManager.Start(this).Test()

	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				err.Error(errors.New("ctx done."))
				return
			case m := <-this.mailbox:
				maybe.TryCatch(
					func() {
						m.Process(this).Test()
					}, nil)
			}
		}
	}(ctx)

	err.Error(nil)
	return
}

func (this *defaultActor) Receive(msg interfaces.Message) (err maybe.MaybeError) {
	if this.mailbox == nil {
		err.Error(errors.New("mailbox not inited."))
		return
	}

	this.mailbox <- msg
	return
}
