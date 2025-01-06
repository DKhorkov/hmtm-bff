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

type InvalidFileExtensionError struct {
	Message string
	BaseErr error
}

func (e InvalidFileExtensionError) Error() string {
	template := "invalid file extension=%s"
	if e.BaseErr != nil {
		return fmt.Sprintf(template+". Base error: %v", e.Message, e.BaseErr)
	}

	return fmt.Sprintf(template, e.Message)
}

type InvalidFileSizeError struct {
	Message string
	BaseErr error
}

func (e InvalidFileSizeError) Error() string {
	template := "invalid file size=%s"
	if e.BaseErr != nil {
		return fmt.Sprintf(template+". Base error: %v", e.Message, e.BaseErr)
	}

	return fmt.Sprintf(template, e.Message)
}
