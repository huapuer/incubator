package actor

import (
	"context"
	"errors"
	"fmt"
	"../common/maybe"
	"../config"
	"../message"
	"strconv"
	"sync"
)

const (
	spikeActorClassName = "actor.spikeActor"
)

func init() {
	RegisterActorPrototype(spikeActorClassName, &defaultActor{}).Test()
}

type spikeActor struct {
	mailbox chan message.Message
	waked   bool
	ctx     context.Context
	mutex   *sync.Mutex
}

func (this spikeActor) New(cfg config.Config, args ...int32) config.IOC {
	ret := MaybeActor{}
	if actor, ok := cfg.Actors[spikeActorClassName]; ok {
		attrs := actor.Attributes.(map[string]string)
		if size, ok := attrs["MailBoxSize"]; ok {
			if mailBoxSize, err := strconv.Atoi(size); err == nil {
				return NewSpikeActor(int64(mailBoxSize))
			}
			ret.Error(fmt.Errorf("illegal actor attribute: %s=%s", "MailBoxSize", size))
			return ret
		}
		ret.Error(fmt.Errorf("no actor attribute found: %s", "MailBoxSize"))
		return ret
	}
	ret.Error(fmt.Errorf("no actor class cfg found: %s", spikeActorClassName))
	return ret
}

func NewSpikeActor(taskChanSize int64) (this MaybeActor) {
	if taskChanSize <= 0 {
		this.Error(fmt.Errorf("wrong task chan size: %d", taskChanSize))
		return
	}
	this.Value(&spikeActor{
		mailbox: make(chan message.Message, taskChanSize),
		mutex:   &sync.Mutex{}})
	return
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
