package context

import "context"

type MessageContext struct{
	context.Context

	Topo int32
}
