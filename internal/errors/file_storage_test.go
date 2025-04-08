package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUploadFileError(t *testing.T) {
	testCases := []struct {
		name           string
		err            UploadFileError
		expectedString string
		expectedBase   error
	}{
		{
			name: "without base error",
			err: UploadFileError{
				Message: "key1",
				BaseErr: nil,
			},
			expectedString: "failed to upload file with key=key1",
			expectedBase:   nil,
		},
		{
			name: "with base error",
			err: UploadFileError{
				Message: "key1",
				BaseErr: errors.New("upload failed"),
			},
			expectedString: "failed to upload file with key=key1. Base error: upload failed",
			expectedBase:   errors.New("upload failed"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualString := tc.err.Error()
			actualBase := tc.err.Unwrap()

			require.Equal(t, tc.expectedString, actualString)
			require.Equal(t, tc.expectedBase, actualBase)
		})
	}
}

func TestInvalidFileExtensionError(t *testing.T) {
	testCases := []struct {
		name           string
		err            InvalidFileExtensionError
		expectedString string
		expectedBase   error
	}{
		{
			name: "without base error",
			err: InvalidFileExtensionError{
				Message: ".exe",
				BaseErr: nil,
			},
			expectedString: "invalid file extension=.exe",
			expectedBase:   nil,
		},
		{
			name: "with base error",
			err: InvalidFileExtensionError{
				Message: ".exe",
				BaseErr: errors.New("extension not allowed"),
			},
			expectedString: "invalid file extension=.exe. Base error: extension not allowed",
			expectedBase:   errors.New("extension not allowed"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualString := tc.err.Error()
			actualBase := tc.err.Unwrap()

			require.Equal(t, tc.expectedString, actualString)
			require.Equal(t, tc.expectedBase, actualBase)
		})
	}
}

func TestInvalidFileSizeError(t *testing.T) {
	testCases := []struct {
		name           string
		err            InvalidFileSizeError
		expectedString string
		expectedBase   error
	}{
		{
			name: "without base error",
			err: InvalidFileSizeError{
				Message: "1024",
				BaseErr: nil,
			},
			expectedString: "invalid file size=1024",
			expectedBase:   nil,
		},
		{
			name: "with base error",
			err: InvalidFileSizeError{
				Message: "1024",
				BaseErr: errors.New("size exceeds limit"),
			},
			expectedString: "invalid file size=1024. Base error: size exceeds limit",
			expectedBase:   errors.New("size exceeds limit"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualString := tc.err.Error()
			actualBase := tc.err.Unwrap()

			require.Equal(t, tc.expectedString, actualString)
			require.Equal(t, tc.expectedBase, actualBase)
		})
	}
}

func TestDeleteFileError(t *testing.T) {
	testCases := []struct {
		name           string
		err            DeleteFileError
		expectedString string
		expectedBase   error
	}{
		{
			name: "without base error",
			err: DeleteFileError{
				Message: "key2",
				BaseErr: nil,
			},
			expectedString: "failed to delete file with key=key2",
			expectedBase:   nil,
		},
		{
			name: "with base error",
			err: DeleteFileError{
				Message: "key2",
				BaseErr: errors.New("delete failed"),
			},
			expectedString: "failed to delete file with key=key2. Base error: delete failed",
			expectedBase:   errors.New("delete failed"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualString := tc.err.Error()
			actualBase := tc.err.Unwrap()

			require.Equal(t, tc.expectedString, actualString)
			require.Equal(t, tc.expectedBase, actualBase)
		})
	}
}
