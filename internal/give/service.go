package give

import (
	"Gives_SDT_Bot/pkg/logging"
	"context"
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
		Owner: ownerId,
		Title: giveTitle,
	}

	if err := s.repository.Create(context.TODO(), give); err != nil {
		s.logger.Errorf("cannot create give with title %s: %v", giveTitle, err)
	}

	return give.Id, nil
}
