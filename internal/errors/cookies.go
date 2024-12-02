package errors

import "fmt"

type CookieNotFoundError struct {
	Message string
}

func (e CookieNotFoundError) Error() string {
	defaultMessage := "cookie not found"
	if e.Message != "" {
		return fmt.Sprintf("\"%s\" %s", e.Message, defaultMessage)
	}

	return defaultMessage
}
