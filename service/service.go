package service

import (
	"context"
	"syncer/errors"
	"syncer/interfaces"
	"time"

	"syncer/utils"
)

type ExternalService struct {
	ProcessLimit uint64
	Duration     time.Duration
}

func (e *ExternalService) GetLimits() (uint64, time.Duration) {

	if utils.GetRandomIntRange(0, 10) == 3 {
		return 0, 0
	}

	return e.ProcessLimit, e.Duration
}

func (e *ExternalService) Process(ctx context.Context, batch interfaces.Batch) error {
	var err error

	if utils.GetRandomIntRange(0, 5) == 3 {
		err = errors.ErrBlocked
	}

	return err
}

type ExternalCoolService struct {
	ProcessLimit uint64
	Duration     time.Duration
}

func (e *ExternalCoolService) GetLimits() (uint64, time.Duration) {
	return e.ProcessLimit, e.Duration
}

func (e *ExternalCoolService) Process(ctx context.Context, batch interfaces.Batch) error {
	return nil
}

type ExternalBadService struct {
	ProcessLimit uint64
	Duration     time.Duration
}

func (e *ExternalBadService) GetLimits() (uint64, time.Duration) {
	return e.ProcessLimit, e.Duration
}

func (e *ExternalBadService) Process(ctx context.Context, batch interfaces.Batch) error {
	return errors.ErrBlocked
}
