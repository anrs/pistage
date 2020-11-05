package executor

import (
	"context"

	"github.com/projecteru2/pistage/action"
)

// Executor .
type Executor interface {
	AsyncStart(ctx context.Context, complex *action.Complex) (jobID string, err error)
	SyncStart(ctx context.Context, complex *action.Complex) (jobID string, err error)
}
