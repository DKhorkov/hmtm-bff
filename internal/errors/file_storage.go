package errors

import "fmt"

type UploadFileError struct {
	Message string
	BaseErr error
}

func (e UploadFileError) Error() string {
	template := "failed to upload file with key=%s"
	if e.BaseErr != nil {
		return fmt.Sprintf(template+". Base error: %v", e.Message, e.BaseErr)
	}

	return fmt.Sprintf(template, e.Message)
}
