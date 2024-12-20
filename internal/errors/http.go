package errors

import "fmt"

type HTTPHandlerTimeoutError struct {
	Message string
	BaseErr error
}

func (e HTTPHandlerTimeoutError) Error() string {
	template := "reached timeout for handling http request"
	if e.Message != "" {
		template = e.Message
	}

	if e.BaseErr != nil {
		return fmt.Sprintf(template+". Base error: %v", e.BaseErr)
	}

	return template
}
