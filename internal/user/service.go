package user

import (
	"Gives_SDT_Bot/pkg/logging"
	"Gives_SDT_Bot/pkg/utils"
	"context"
	"fmt"
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
		TgID: utils.Int64ToString(telegramId),
	}

	if isAdmin {
		user.IsAdmin = true
	}

	if err := s.repository.Create(context.TODO(), user); err != nil {
		s.logger.Errorf("cannot create user with tgId=%s:\n%s", user.TgID, err)
		return 0, err
	}
	return user.ID, nil
}

func (s *Service) GetUserIdByTgId(telegramId int64) (int, error) {
	user := &User{}

	if err := s.repository.FindOneWithConditions(context.TODO(), user, fmt.Sprintf("tg_id='%d'", telegramId)); err != nil {
		s.logger.Errorf("cannot find user with tgId=%d:\n%s", telegramId, err)
		return 0, err
	}

	return user.ID, nil
}

func (s *Service) GetAdmins() ([]int64, error) {
	var adminsIds []int64

	admins, err := s.repository.FindAllWithConditions(context.TODO(), "is_admin=true")
	if err != nil {
		s.logger.Error("cannot get admins")
		return nil, err
	}

	for _, admin := range admins {
		adminId, err := utils.StringToInt64(admin.TgID)
		if err != nil {
			s.logger.Errorf("cannot parse admin tgId=%s to int64", admin.TgID)
			return nil, err
		}
		adminsIds = append(adminsIds, adminId)
	}

	return adminsIds, nil
}
