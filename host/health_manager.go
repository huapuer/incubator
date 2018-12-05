package host

import (
	"../common/maybe"
	"fmt"
	"github.com/incubator/layer"
	"github.com/incubator/message"
	"time"
)

type HealthManager interface {
	IsHealth() bool
	Heartbeat()
	Health()
	Start() maybe.MaybeError
	GetLayer() int32
	SetLayer(int32)
}

type defaultHealthManager struct {
	checkIntvl     time.Duration
	heartbeatIntvl time.Duration
	lastHealthTime time.Time
	health         bool
	layer          int32
}

func (this *defaultHealthManager) Heartbeat() {
	if time.Now().Sub(this.lastHealthTime) > this.heartbeatIntvl {
		this.health = false
		return
	}
	this.health = true
}

func (this *defaultHealthManager) Health() {
	this.lastHealthTime = time.Now()
}

func (this *defaultHealthManager) Start() (err maybe.MaybeError) {
	if this.checkIntvl <= 0 {
		err.Error(fmt.Errorf("illegal check interval: %d", this.checkIntvl))
		return
	}

	if this.heartbeatIntvl <= 0 {
		err.Error(fmt.Errorf("illegal heartbeat interval: %d", this.heartbeatIntvl))
		return
	}

	if this.checkIntvl > this.heartbeatIntvl {
		err.Error(fmt.Errorf("check interval > heartbeat interval: %d>%d",
			this.checkIntvl, this.heartbeatIntvl))
		return
	}

	this.lastHealthTime = time.Now()

	go func() {
		l := layer.GetLayer(this.layer).Right()
		for {
			for _, h := range l.GetTopo().GetRemoteHosts() {
				msg := &message.HealthCheckReqMessage{}
				msg.SetLayer(int8(this.layer))
				msg.SetSrcLayer(int8(this.layer))
				msg.SetType(int8(l.GetMessageType(msg).Right()))

				t := l.GetTopo()
				msg.SetSrcHostId(t.GetRemoteHostId(int32(h.GetId())))
				t.SendToHost(h.GetId(), msg).Test()
			}
			time.Sleep(this.checkIntvl)
		}
	}()

	go func() {
		for {
			this.Heartbeat()
			time.Sleep(this.heartbeatIntvl)
		}
	}()

	err.Error(nil)
	return
}

func (this defaultHealthManager) GetLayer() int32 {
	return this.layer
}

func (this *defaultHealthManager) SetLayer(layer int32) {
	this.layer = layer
}
