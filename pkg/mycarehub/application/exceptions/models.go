package exceptions

import "fmt"

// CustomError represents a custom error struct
// Reference https://blog.golang.org/error-handling-and-go
type CustomError struct {
	Err     error                  `json:"error,omitempty"`
	Message string                 `json:"message,omitempty"`
	Code    int                    `json:"code,omitempty"`
	Meta    map[string]interface{} `json:"meta"`
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("%d: %s: %v: %v", e.Code, e.Message, e.Err, e.Meta)
}
