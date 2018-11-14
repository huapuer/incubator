package actor

import (
	"context"
	"errors"
	"../common/maybe"
	"../config"
	"../message"
	"sync"
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

	waked   bool
	ctx     context.Context
	mutex   *sync.Mutex
}

func (this spikeActor) New(attrs interface{}, cfg config.Config) config.IOC {
	ret := MaybeActor{}

	actor := &defaultActor{}
	actor.mailBox.Init(attrs, cfg).Test()

	ret.Value(actor)
	return ret
}

func (this *spikeActor) Start(ctx context.Context) (err maybe.MaybeError) {
	if this.mailbox == nil {
		err.Error(errors.New("mailbox not inited."))
		return
	}

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
							m.Process(this.ctx).Test()
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
