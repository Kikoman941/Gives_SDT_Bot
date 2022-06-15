package fsm

import (
	"Gives_SDT_Bot/pkg/logging"
	"context"
)

func NewFSM(repository Repository, logger *logging.Logger) (*FSM, error) {
	return &FSM{
		repository: repository,
		logger:     logger,
	}, nil
}

func (fsm *FSM) SetState(ctx context.Context, us UserState) error {
	if err := fsm.repository.UpdateOrInsert(ctx, us); err != nil {
		return err
	}
	return nil
}
