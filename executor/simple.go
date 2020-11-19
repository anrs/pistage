package executor

import (
	"context"

	"github.com/projecteru2/pistage/action"
	"github.com/projecteru2/pistage/errors"
	"github.com/projecteru2/pistage/orch"
)

// Simple executor.
type Simple struct {
	orch  orch.Orchestrator
	store Store
}

// NewSimple .
func NewSimple() (simple *Simple, err error) {
	simple = &Simple{store: &localStore{}}
	simple.orch, err = orch.NewEru()
	return
}

// AsyncStart .
func (s *Simple) AsyncStart(ctx context.Context, complex *action.Complex) (string, error) {
	return s.start(ctx, complex, true)
}

// SyncStart .
func (s *Simple) SyncStart(ctx context.Context, complex *action.Complex) (string, error) {
	return s.start(ctx, complex, false)
}

func (s *Simple) start(ctx context.Context, complex *action.Complex, async bool) (string, error) {
	jg, err := s.parse(complex)
	if err != nil {
		return "", errors.Trace(err)
	}

	if err := jg.save(ctx); err != nil {
		return "", errors.Trace(err)
	}

	if async {
		// TODO
		// actually, we should booting a labmda as a workload to waiting
		// the really executor has been done.
		go jg.run(ctx)
	} else {
		jg.run(ctx)
	}

	return jg.id, nil
}

func (s *Simple) parse(complex *action.Complex) (*JobGroup, error) {
	jg := NewJobGroup()
	jg.params = complex.Params
	jg.orch = s.orch
	jg.store = s.store

	for name, acts := range complex.Groups {
		if err := jg.add(name, acts); err != nil {
			return jg, errors.Trace(err)
		}
	}

	return jg, nil
}
