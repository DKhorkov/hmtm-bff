package services

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	mocklogging "github.com/DKhorkov/libs/logging/mocks"

	customerrors "github.com/DKhorkov/hmtm-bff/internal/errors"
	mockrepositories "github.com/DKhorkov/hmtm-bff/mocks/repositories"
)

func TestFileStorageService_Upload(t *testing.T) {
	ctrl := gomock.NewController(t)
	fileStorageRepository := mockrepositories.NewMockFileStorageRepository(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	service := NewFileStorageService(fileStorageRepository, logger)

	testCases := []struct {
		name          string
		key           string
		data          []byte
		setupMocks    func(fileStorageRepository *mockrepositories.MockFileStorageRepository, logger *mocklogging.MockLogger)
		expectedURL   string
		errorExpected bool
	}{
		{
			name: "success",
			key:  "test-key",
			data: []byte("test-data"),
			setupMocks: func(fileStorageRepository *mockrepositories.MockFileStorageRepository, logger *mocklogging.MockLogger) {
				fileStorageRepository.
					EXPECT().
					Upload(gomock.Any(), "test-key", []byte("test-data")).
					Return("http://storage/test-key", nil).
					Times(1)
			},
			expectedURL:   "http://storage/test-key",
			errorExpected: false,
		},
		{
			name: "error",
			key:  "test-key",
			data: []byte("test-data"),
			setupMocks: func(fileStorageRepository *mockrepositories.MockFileStorageRepository, logger *mocklogging.MockLogger) {
				fileStorageRepository.
					EXPECT().
					Upload(gomock.Any(), "test-key", []byte("test-data")).
					Return("", errors.New("upload failed")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			expectedURL:   "",
			errorExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(fileStorageRepository, logger)
			}

			url, err := service.Upload(context.Background(), tc.key, tc.data)
			if tc.errorExpected {
				require.Error(t, err)
				require.IsType(t, &customerrors.UploadFileError{}, err)
				require.Equal(t, tc.key, err.(*customerrors.UploadFileError).Message)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.expectedURL, url)
		})
	}
}

func TestFileStorageService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	fileStorageRepository := mockrepositories.NewMockFileStorageRepository(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	service := NewFileStorageService(fileStorageRepository, logger)

	testCases := []struct {
		name          string
		key           string
		setupMocks    func(fileStorageRepository *mockrepositories.MockFileStorageRepository, logger *mocklogging.MockLogger)
		errorExpected bool
	}{
		{
			name: "success",
			key:  "test-key",
			setupMocks: func(fileStorageRepository *mockrepositories.MockFileStorageRepository, logger *mocklogging.MockLogger) {
				fileStorageRepository.
					EXPECT().
					Delete(gomock.Any(), "test-key").
					Return(nil).
					Times(1)
			},
			errorExpected: false,
		},
		{
			name: "error",
			key:  "test-key",
			setupMocks: func(fileStorageRepository *mockrepositories.MockFileStorageRepository, logger *mocklogging.MockLogger) {
				fileStorageRepository.
					EXPECT().
					Delete(gomock.Any(), "test-key").
					Return(errors.New("delete failed")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			errorExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(fileStorageRepository, logger)
			}

			err := service.Delete(context.Background(), tc.key)
			if tc.errorExpected {
				require.Error(t, err)
				require.Equal(t, "delete failed", err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestFileStorageService_DeleteMany(t *testing.T) {
	ctrl := gomock.NewController(t)
	fileStorageRepository := mockrepositories.NewMockFileStorageRepository(ctrl)
	logger := mocklogging.NewMockLogger(ctrl)
	service := NewFileStorageService(fileStorageRepository, logger)

	testCases := []struct {
		name           string
		keys           []string
		setupMocks     func(fileStorageRepository *mockrepositories.MockFileStorageRepository, logger *mocklogging.MockLogger)
		expectedErrors []error
	}{
		{
			name: "success",
			keys: []string{"key1", "key2"},
			setupMocks: func(fileStorageRepository *mockrepositories.MockFileStorageRepository, logger *mocklogging.MockLogger) {
				fileStorageRepository.
					EXPECT().
					DeleteMany(gomock.Any(), []string{"key1", "key2"}).
					Return(nil).
					Times(1)
			},
			expectedErrors: nil,
		},
		{
			name: "partial errors",
			keys: []string{"key1", "key2"},
			setupMocks: func(fileStorageRepository *mockrepositories.MockFileStorageRepository, logger *mocklogging.MockLogger) {
				deleteErrors := []error{nil, errors.New("delete key2 failed")}
				fileStorageRepository.
					EXPECT().
					DeleteMany(gomock.Any(), []string{"key1", "key2"}).
					Return(deleteErrors).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			expectedErrors: []error{nil, errors.New("delete key2 failed")},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(fileStorageRepository, logger)
			}

			deleteErrors := service.DeleteMany(context.Background(), tc.keys)
			if tc.expectedErrors == nil {
				require.Nil(t, deleteErrors)
			} else {
				require.Equal(t, len(tc.expectedErrors), len(deleteErrors))
				for i, expectedErr := range tc.expectedErrors {
					if expectedErr == nil {
						require.Nil(t, deleteErrors[i])
					} else {
						require.Equal(t, expectedErr.Error(), deleteErrors[i].Error())
					}
				}
			}
		})
	}
}
