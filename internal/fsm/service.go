package fsm

import (
	"Gives_SDT_Bot/internal/fsm/db"
	"Gives_SDT_Bot/pkg/client/postgresql"
	"Gives_SDT_Bot/pkg/logging"
)

type Service struct {
	repository Repository
	logger     *logging.Logger
}

func NewFSMService(dbClient postgresql.Client, logger *logging.Logger) *Service {
	return &Service{
		repository: db.NewRepository(dbClient, logger),
		logger:     logger,
	}
}
