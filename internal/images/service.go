package images

import (
	"Gives_SDT_Bot/pkg/logging"
	"fmt"
	"gopkg.in/telebot.v3"
	"path/filepath"
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

func (s *Service) SaveFile(file *telebot.File, giveId string) (string, error) {
	filePath := fmt.Sprintf("./.images/%s.jpeg", giveId)
	filename := filepath.Base(filePath)

	if err := s.repository.Download(file, filePath); err != nil {
		s.logger.Errorf("caannot download fileId=%s: %s", file.FileID, err)
		return "", err
	}
	return filename, nil
}
