package usecases

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateFileExtension(t *testing.T) {
	testCases := []struct {
		name              string
		extension         string
		allowedExtensions []string
		expected          bool
	}{
		{
			name:              "Valid extension",
			extension:         ".png",
			allowedExtensions: []string{".png", ".jpg", ".jpeg"},
			expected:          true,
		},
		{
			name:              "Invalid extension",
			extension:         ".txt",
			allowedExtensions: []string{".png", ".jpg", ".jpeg"},
			expected:          false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			validExtension := validateFileExtension(tc.extension, tc.allowedExtensions)
			assert.Equal(t, tc.expected, validExtension)
		})
	}
}

func TestValidateFileSize(t *testing.T) {
	testCases := []struct {
		name     string
		size     int64
		maxSize  int64
		expected bool
	}{
		{
			name:     "Valid size",
			size:     1024,
			maxSize:  1024,
			expected: true,
		},
		{
			name:     "Invalid size",
			size:     2048,
			maxSize:  1024,
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			validSize := validateFileSize(tc.size, tc.maxSize)
			assert.Equal(t, tc.expected, validSize)
		})
	}
}
