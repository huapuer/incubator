package context

import (
	"context"
	"github.com/incubator/common/maybe"
)

type AyncErrorContext struct {
	Ctx context.Context
	Err chan maybe.MaybeError
}
