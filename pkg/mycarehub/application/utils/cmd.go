package utils

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseChoice validates a choice is in the range of a choices
func ParseChoice[T any](choices []T, choice string) (*T, error) {
	choice = strings.TrimSpace(choice)
	index, err := strconv.Atoi(choice)
	if err != nil {
		return nil, err
	}

	if index >= len(choices) || index < 0 {
		return nil, fmt.Errorf("invalid choice: %v", choice)
	}

	return &choices[index], nil
}
