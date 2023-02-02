package utils

import (
	"io"
	"os"
	"path/filepath"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
)

// ReadFile reads the content of a file and return a slice of bytes
func ReadFile(path string) ([]byte, error) {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(absolutePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := file.Close()
		if err != nil {
			helpers.ReportErrorToSentry(err)
		}
	}()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return byteValue, nil
}
