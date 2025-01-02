package errors

import "fmt"

type PermissionDeniedError struct {
	Message string
	BaseErr error
}

func (e PermissionDeniedError) Error() string {
	template := "permission denied"
	if e.Message != "" {
		template = fmt.Sprintf(template+": %s", e.Message)
	}

	if e.BaseErr != nil {
		return fmt.Sprintf(template+". Base error: %v", e.BaseErr)
	}

	return template
}
