package context

import (
	"context"
	"../host"
)

type HostRecoverContext struct{
	context.Context

	Ctx context.Context
	Ret chan host.MaybeHost
}

func NewHostRecoverContext() (ret HostRecoverContext, cancel context.CancelFunc){
	ctx, cancel := context.WithCancel(context.Background())
	return HostRecoverContext{
		Ctx: ctx,
		Ret: make(chan host.MaybeHost),
	}, cancel
}
