package customError

import (
	"fmt"
)

type CustomError struct {
	Code    string
	Message string
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("Code: %s, Message: %s", e.Code, e.Message)
}

func New(code string, message string) *CustomError {
	return &CustomError{
		Code:    code,
		Message: message,
	}
}
