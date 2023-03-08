package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
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

// LoadScreeningTools is a helper function that loads the given screening tools when a program is created
// TODO: temporary solution. this should be removed after implementing screening tool creation in the frontend
func LoadScreeningTools() ([]*dto.ScreeningToolInput, error) {
	screeningTools := []*dto.ScreeningToolInput{}
	bs, err := ReadFile("data/mycarehub_program_screeningtools.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read screening tool file:%w", err)
	}
	err = json.Unmarshal(bs, &screeningTools)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal screening tool:%w", err)
	}
	return screeningTools, nil
}
