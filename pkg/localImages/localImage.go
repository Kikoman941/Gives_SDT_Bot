package localImages

import (
	"Gives_SDT_Bot/pkg/errors"
	"Gives_SDT_Bot/pkg/logging"
	"os"
	"path/filepath"
)

type Client interface {
	Save(img string) error
}

type LocalImages struct {
	place  string // Folder will be created at root of project
	logger *logging.Logger
}

func NewLocalImage(place string, logger *logging.Logger) (*LocalImages, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, errors.FormatError("cannot get current working directory", err)
	}

	path := filepath.Join(cwd, place)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.Mkdir(path, os.ModePerm); err != nil {
			return nil, errors.FormatError("cannot create directory", err)
		}
	} else {
		logger.Info("Directory for images already exists, used this: " + path)
	}

	return &LocalImages{
		place:  path,
		logger: logger,
	}, nil
}

func (li *LocalImages) Save(img string) error {
	return nil
}
