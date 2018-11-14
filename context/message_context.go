package context

import (
	"context"
	"github.com/incubator/actor"
)

type MessageContext struct{
	context.Context

	Runner actor.Actor
}
