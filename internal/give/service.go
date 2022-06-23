package give

import (
	"Gives_SDT_Bot/pkg/logging"
	"context"
	"github.com/davecgh/go-spew/spew"
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
	gives, err := s.repository.FindAllWithConditions(context.TODO(), spew.Sprintf("owner=%d", userId))
	if err != nil {
		s.logger.Errorf("cannot get userId=%d gives: %s", userId, err)
	}

	return gives, nil
}
