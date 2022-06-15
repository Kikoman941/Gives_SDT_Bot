package errors

import (
	"errors"
	"fmt"
)

func FormatError(text string, error error) error {
	return errors.New(
		fmt.Sprintf("%s:\n%s", text, error),
	)
}
