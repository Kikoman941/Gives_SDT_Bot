package images

import (
	"Gives_SDT_Bot/pkg/logging"
)

type Service struct {
	repository Repository
	logger     *logging.Logger
}

func NewImagesService(repository Repository, logger *logging.Logger) *Service {
	return &Service{
		repository: repository,
		logger:     logger,
	}
}
