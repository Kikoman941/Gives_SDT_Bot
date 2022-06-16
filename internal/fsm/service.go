package fsm

import (
	"Gives_SDT_Bot/internal/fsm/db"
	"Gives_SDT_Bot/pkg/client/postgresql"
	"Gives_SDT_Bot/pkg/logging"
	"context"
)

type Service struct {
	repository Repository
	logger     *logging.Logger
}

func NewFSMService(dbClient postgresql.Client, logger *logging.Logger) *Service {
	return &Service{
		repository: db.NewRepository(dbClient, logger),
		logger:     logger,
	}
}

func (s *Service) SetState(userId int, state string) error {
	userState := &UserState{
		UserID: userId,
		State:  state,
	}
	err := s.repository.UpdateOrInsert(context.TODO(), userState)
	if err != nil {
		s.logger.Errorf("cannot set state=%s for user with tgId=%d", state, userId)
		return err
	}
	return nil
}
