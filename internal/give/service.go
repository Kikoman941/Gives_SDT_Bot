package give

import (
	"Gives_SDT_Bot/pkg/logging"
	"context"
	"fmt"
)

type Service struct {
	repository Repository
	logger     *logging.Logger
}

func NewGiveService(repository Repository, logger *logging.Logger) *Service {
	return &Service{
		repository: repository,
		logger:     logger,
	}
}

func (s *Service) CreateGive(giveTitle string, ownerId int) (int, error) {
	give := &Give{
		IsActive: false,
		Owner:    ownerId,
		Title:    giveTitle,
	}

	if err := s.repository.Create(context.TODO(), give); err != nil {
		s.logger.Errorf("cannot create give with title %s: %s", giveTitle, err)
	}

	return give.Id, nil
}

func (s *Service) GetAllUserGives(userId int) ([]Give, error) {
	gives, err := s.repository.FindAllWithConditions(context.TODO(), fmt.Sprintf("owner=%d", userId))
	if err != nil {
		s.logger.Errorf("cannot get userId=%d gives: %s", userId, err)
	}

	return gives, nil
}

func (s *Service) UpdateGive(giveId int, update string) error {
	err := s.repository.UpdateWithConditions(
		context.TODO(),
		fmt.Sprintf("id=%d", giveId),
		update,
	)
	if err != nil {
		s.logger.Errorf("cannot do update giveId=%d: %s", giveId, err)
		return err
	}

	return nil
}
