package images

import (
	"Gives_SDT_Bot/internal/images/db"
	"Gives_SDT_Bot/pkg/localImages"
	"Gives_SDT_Bot/pkg/logging"
)

type Service struct {
	repository Repository
	logger     *logging.Logger
}

func NewImagesService(client localImages.Client, logger *logging.Logger) *Service {
	return &Service{
		repository: db.NewRepository(client, logger),
		logger:     logger,
	}
}
