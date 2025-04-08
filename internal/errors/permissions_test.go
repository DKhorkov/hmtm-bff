package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPermissionDeniedError(t *testing.T) {
	testCases := []struct {
		name           string
		err            PermissionDeniedError
		expectedString string
		expectedBase   error
	}{
		{
			name: "default message without base error",
			err: PermissionDeniedError{
				Message: "",
				BaseErr: nil,
			},
			expectedString: "permission denied",
			expectedBase:   nil,
		},
		{
			name: "default message with base error",
			err: PermissionDeniedError{
				Message: "",
				BaseErr: errors.New("access restricted"),
			},
			expectedString: "permission denied. Base error: access restricted",
			expectedBase:   errors.New("access restricted"),
		},
		{
			name: "custom message without base error",
			err: PermissionDeniedError{
				Message: "user not authorized",
				BaseErr: nil,
			},
			expectedString: "permission denied: user not authorized",
			expectedBase:   nil,
		},
		{
			name: "custom message with base error",
			err: PermissionDeniedError{
				Message: "user not authorized",
				BaseErr: errors.New("access restricted"),
			},
			expectedString: "permission denied: user not authorized. Base error: access restricted",
			expectedBase:   errors.New("access restricted"),
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
