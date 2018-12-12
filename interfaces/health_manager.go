package interfaces

import "github.com/incubator/common/maybe"

type HealthManager interface {
	IsHealth() bool
	Heartbeat()
	Health()
	Start() maybe.MaybeError
	GetLayer() int32
	SetLayer(int32)
}
