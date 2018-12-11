package actor

import (
	"../common/maybe"
	"../config"
	"../message"
	"context"
	"errors"
	"sync"
	"time"
)

const (
	spikeActorClassName = "actor.spikeActor"
)

func init() {
	RegisterActorPrototype(spikeActorClassName, &defaultActor{}).Test()
}

type spikeActor struct {
	commonActor
	mailBox
	defaultHealthManager

	waked bool
	ctx   context.Context
	mutex *sync.Mutex
}

func (this spikeActor) New(attrs interface{}, cfg config.Config) config.IOC {
	ret := MaybeActor{}

	heartbeatIntvl := config.GetAttrInt64(attrs, "HeartbeatIntvl", config.CheckInt64GT0).Right()

	actor := &spikeActor{
		defaultHealthManager: defaultHealthManager{
			heartbeatIntvl: time.Duration(heartbeatIntvl),
		},
	}
	actor.mailBox.Init(attrs, cfg).Test()

	ret.Value(actor)
	return ret
}

func (this *spikeActor) Start(ctx context.Context) (err maybe.MaybeError) {
	if this.mailbox == nil {
		err.Error(errors.New("mailbox not inited."))
		return
	}

	this.defaultHealthManager.Start(this).Test()

	this.ctx = ctx

	return
}

func (this *spikeActor) Receive(msg message.Message) (err maybe.MaybeError) {
	if this.mailbox == nil {
		err.Error(errors.New("mailbox not inited."))
		return
	}

	this.mailbox <- msg

	this.mutex.Lock()
	if !this.waked {
		go func() {
			for {
				processed := false
				select {
				case <-this.ctx.Done():
					err.Error(errors.New("ctx done."))
					return
				case m := <-this.mailbox:
					maybe.TryCatch(
						func() {
							m.Process(this).Test()
						}, nil)
					processed = true
				}
				if !processed {
					break
				}
			}
			this.waked = false
		}()

		this.waked = true
	}
	this.mutex.Unlock()

	return
}
