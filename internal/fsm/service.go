package fsm

import (
	"Gives_SDT_Bot/pkg/logging"
	"context"
	"fmt"
)

type Service struct {
	repository Repository
	logger     *logging.Logger
}

func NewFSMService(repository Repository, logger *logging.Logger) *Service {
	return &Service{
		repository: repository,
		logger:     logger,
	}
}

func (s *Service) SetState(userId int, state string) error {
	userState := &UserState{
		UserID: userId,
		State:  state,
	}
	err := s.repository.InsertOrUpdate(context.TODO(), userState)
	if err != nil {
		s.logger.Errorf("cannot set state=%s for user with tgId=%d:\n%s", state, userId, err)
		return err
	}
	return nil
}

func (s *Service) GetState(userId int) (string, error) {
	userState := &UserState{}
	err := s.repository.FindOneWithConditions(context.TODO(), userState, fmt.Sprintf("user_id=%d", userId))
	if err != nil {
		s.logger.Errorf("cannot get userId=%d state:\n%s", userId, err)
		return "", err
	}

	return userState.State, nil
}
