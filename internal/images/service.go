package images

import (
	"Gives_SDT_Bot/pkg/logging"
	"fmt"
	"gopkg.in/telebot.v3"
	"os"
	"path/filepath"
)

type Service struct {
	repository Repository
	logger     *logging.Logger
}

func NewImagesService(repository Repository, logger *logging.Logger) *Service {
	cwd, err := os.Getwd()
	if err != nil {
		logger.Fatalf("cannot get current working directory: %s", err)
		return nil
	}

	path := filepath.Join(cwd, ".images")

	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.Mkdir(path, os.ModePerm); err != nil {
			logger.Fatalf("cannot create directory: %s", err)
			return nil
		}
	} else {
		logger.Info("Directory for images already exists, used this: " + path)
	}

	return &Service{
		repository: repository,
		logger:     logger,
	}
}

func (s *Service) SaveFile(file *telebot.File, giveId string) (string, error) {
	filePath := fmt.Sprintf("./.images/%s.jpeg", giveId)
	filename := filepath.Base(filePath)

	if err := s.repository.Download(file, filePath); err != nil {
		s.logger.Errorf("cannot download fileId=%s: %s", file.FileID, err)
		return "", err
	}
	return filename, nil
}
