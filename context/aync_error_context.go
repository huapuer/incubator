package context

import (
	"../common/maybe"
	"context"
)

type AyncErrorContext struct {
	Ctx context.Context
	Err chan maybe.MaybeError
}
