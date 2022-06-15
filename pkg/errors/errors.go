package errors

import (
	"fmt"
)

func FormatError(text string, error error) error {
	return fmt.Errorf("%s:\n%s", text, error)
}
