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

func (this spikeActor) New(attrs interface{}, cfg config.Config) config.IOC {
	ret := MaybeActor{}
	if attrs == nil{
		ret.Error(errors.New("actor attrs is nil"))
		return ret
	}
	attrsMap, ok := attrs.(map[string]interface{})
	if !ok{
		ret.Error(fmt.Errorf("illegal cfg type when new actor %s", spikeActorClassName))
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
