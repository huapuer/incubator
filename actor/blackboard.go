package actor

import (
	"errors"
	"fmt"
	"github.com/incubator/common/maybe"
	"github.com/incubator/interfaces"
	"github.com/incubator/message"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

type blackBoard struct {
	data  map[string]interface{}
	token map[string]int64
}

func (this *blackBoard) SetState(runner interfaces.Actor, key string, value interface{}, expire time.Duration, expireFunc func(interfaces.Actor)) (err maybe.MaybeError) {
	if this.data == nil {
		this.data = make(map[string]interface{})
	}
	if _, ok := this.data[key]; ok {
		err.Error(fmt.Errorf("state key already exists: %s", key))
		return
	}
	this.data[key] = value
	token := rand.Int63()
	this.data[key] = token
	if expire > 0 {
		go func() {
			<-time.After(expire)
			runner.Receive(message.StateExpireMessage{
				key,
				token,
				expireFunc})
		}()
	}
	err.Error(nil)
	return
}

func (this *blackBoard) UnsetState(key string) (err maybe.MaybeError) {
	if this.data == nil {
		err.Error(errors.New("actor blackboard not set"))
		return
	}
	delete(this.data, key)
	err.Error(nil)
	return
}

func (this *blackBoard) UnsetStateWithToken(key string, token int64) (err maybe.MaybeError) {
	if this.data == nil {
		err.Error(errors.New("actor blackboard not set"))
		return
	}
	t, ok := this.token[key]
	if !ok {
		err.Error(fmt.Errorf("token not exists: %d", token))
		return
	}
	if t != token {
		err.Error(errors.New("token not match"))
		return
	}
	delete(this.data, key)
	err.Error(nil)
	return
}

func (this blackBoard) GetState(key string) (ret maybe.MaybeEface) {
	if this.data == nil {
		ret.Error(errors.New("actor blackboard not set"))
		return
	}
	v, ok := this.data[key]
	if !ok {
		ret.Error(fmt.Errorf("state key does not exists: %s", key))
		return
	}
	ret.Value(v)
	return
}
