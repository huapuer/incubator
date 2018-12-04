package host

import (
	"time"
	"incubator/common/maybe"
	"fmt"
)

type HealthManager interface{
	IsHealth() bool
	Heartbeat()
	Health()
	Start() maybe.MaybeError
}

type defaultHealthManager struct {
	heartbeatIntvl time.Duration
	lastHealthTime time.Time
	health bool
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
	if this.heartbeatIntvl <= 0 {
		err.Error(fmt.Errorf("illegal heartbeat interval: %d", this.heartbeatIntvl))
		return
	}

	this.lastHealthTime = time.Now()

	//TODO: send health_check_req_message periodically

	go func() {
		for {
			this.Heartbeat()
			time.Sleep(this.heartbeatIntvl)
		}
	}()

	err.Error(nil)
	return
}
