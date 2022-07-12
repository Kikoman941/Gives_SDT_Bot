package user

import (
	"Gives_SDT_Bot/pkg/logging"
	"Gives_SDT_Bot/pkg/utils"
	"context"
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
	users, err := s.repository.FindAllWithConditions(context.TODO(), `"tgId"=?`, utils.Int64ToString(telegramId))
	if err != nil {
		s.logger.Errorf("cannot find user with tgId=%d: %s", telegramId, err)
		return 0, err
	}

	return users[0].ID, nil
}

func (s *Service) GetTgIdByUserId(userId int) (string, error) {
	users, err := s.repository.FindAllWithConditions(context.TODO(), `"id"=?`, userId)
	if err != nil {
		s.logger.Errorf("cannot find user userId=%d: %s", userId, err)
		return "", err
	}

	return users[0].TgID, nil
}

func (s *Service) GetAdmins() ([]int64, error) {
	var adminsIds []int64

	admins, err := s.repository.FindAllWithConditions(context.TODO(), `"isAdmin"=?`, true)
	if err != nil {
		s.logger.Errorf("cannot get admins: %s", err)
		return nil, err
	}

	for _, admin := range admins {
		adminId, err := utils.StringToInt64(admin.TgID)
		if err != nil {
			s.logger.Errorf("cannot parse admin tgId=%s to int64: %s", admin.TgID, err)
			return nil, err
		}
		adminsIds = append(adminsIds, adminId)
	}

	return adminsIds, nil
}
