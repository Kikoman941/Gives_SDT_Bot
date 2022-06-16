package db

import (
	"Gives_SDT_Bot/internal/images"
	"Gives_SDT_Bot/pkg/localImages"
	"Gives_SDT_Bot/pkg/logging"
)

type repository struct {
	client localImages.Client
	logger *logging.Logger
}

func (r *repository) SaveImage(img string) error {
	//TODO implement me
	panic("implement me")
}

func NewRepository(client localImages.Client, logger *logging.Logger) images.Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}
