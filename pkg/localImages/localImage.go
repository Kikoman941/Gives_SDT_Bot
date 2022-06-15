package localImages

import (
	"Gives_SDT_Bot/pkg/errors"
	"fmt"
	"os"
	"path/filepath"
)

type LocalImages struct {
	place string // Folder will be created at root of project
}

func NewLocalImage(place string) (*LocalImages, error) {
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
		fmt.Println("Directory for images already exists, used this: " + path)
	}

	return &LocalImages{
		place: path,
	}, nil
}
