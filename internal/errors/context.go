package errors

import "fmt"

type ContextValueNotFoundError struct {
	Message string
}

func (e ContextValueNotFoundError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("context with value \"%s\" not found", e.Message)
	}

	return "context with value not found"
}
