package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHTTPHandlerTimeoutError(t *testing.T) {
	testCases := []struct {
		name           string
		err            HTTPHandlerTimeoutError
		expectedString string
		expectedBase   error
	}{
		{
			name: "default message without base error",
			err: HTTPHandlerTimeoutError{
				Message: "",
				BaseErr: nil,
			},
			expectedString: "reached timeout for handling http request",
			expectedBase:   nil,
		},
		{
			name: "default message with base error",
			err: HTTPHandlerTimeoutError{
				Message: "",
				BaseErr: errors.New("context deadline exceeded"),
			},
			expectedString: "reached timeout for handling http request. Base error: context deadline exceeded",
			expectedBase:   errors.New("context deadline exceeded"),
		},
		{
			name: "custom message without base error",
			err: HTTPHandlerTimeoutError{
				Message: "custom timeout error",
				BaseErr: nil,
			},
			expectedString: "custom timeout error",
			expectedBase:   nil,
		},
		{
			name: "custom message with base error",
			err: HTTPHandlerTimeoutError{
				Message: "custom timeout error",
				BaseErr: errors.New("context deadline exceeded"),
			},
			expectedString: "custom timeout error. Base error: context deadline exceeded",
			expectedBase:   errors.New("context deadline exceeded"),
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
