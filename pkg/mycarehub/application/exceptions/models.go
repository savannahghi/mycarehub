package exceptions

import "fmt"

// CustomError represents a custom error struct
// Reference https://blog.golang.org/error-handling-and-go
type CustomError struct {
	Err     error  `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code,omitempty"`
	Detail  string `json:"detail,omitempty"`
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("%d: %s: %s", e.Code, e.Err, e.Detail)
}

// GetErrorCode returns the error code from custom error
func GetErrorCode(err error) int {
	if err == nil {
		return int(Internal)
	}
	if e, ok := err.(*CustomError); ok {
		return e.Code
	}
	return int(Internal)
}

// GetError returns the error from custom error
func GetError(err error) error {
	if err == nil {
		return nil
	}
	if e, ok := err.(*CustomError); ok {
		return e.Err
	}
	return err
}
