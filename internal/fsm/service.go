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

func (s *Service) Setstate(userId int, state string, data map[string]string) error {
	userstate := &Userstate{
		UserID: userId,
		state:  state,
		Data:   data,
	}
	err := s.repository.InsertOrUpdate(context.TODO(), userstate)
	if err != nil {
		s.logger.Errorf("cannot set (state=%s data=%s) for user with tgId=%d: %s", state, data, userId, err)
		return err
	}
	return nil
}

func (s *Service) Getstate(userId int) (*Userstate, error) {
	usersstates, err := s.repository.FindAllWithConditions(context.TODO(), fmt.Sprintf("\"userId\"=%d", userId))
	if err != nil {
		s.logger.Errorf("cannot get userId=%d state: %s", userId, err)
		return nil, err
	}

	return &usersstates[0], nil
}

func NewFSMService(repository Repository, logger *logging.Logger) *Service {
	return &Service{
		repository: repository,
		logger:     logger,
	}
}
