package context

import (
	"context"
	"../common/maybe"
)

type AyncErrorContext struct {
	Ctx context.Context
	Err chan maybe.MaybeError
}
