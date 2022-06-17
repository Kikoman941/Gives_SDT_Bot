package user

import (
	"Gives_SDT_Bot/pkg/logging"
	"context"
	"strconv"
)

type Service struct {
	repository Repository
	logger     *logging.Logger
}

func NewUserService(repository Repository, logger *logging.Logger) *Service {
	return &Service{
		repository: repository,
		logger:     logger,
	}
}

func (s *Service) AddUser(telegramId int64, isAdmin bool) (int, error) {
	user := &User{
		TgID: strconv.FormatInt(telegramId, 10),
	}

	if isAdmin {
		user.IsAdmin = true
	}

	userId, err := s.repository.Create(context.TODO(), user)
	if err != nil {
		s.logger.Errorf("cannot create user with tgId=%s:\n%s", user.TgID, err)
		return 0, err
	}
	return userId, nil
}
