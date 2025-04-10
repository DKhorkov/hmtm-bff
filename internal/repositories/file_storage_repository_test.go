package repositories

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	mocklogging "github.com/DKhorkov/libs/logging/mocks"
	"github.com/DKhorkov/libs/pointers"

	appconfig "github.com/DKhorkov/hmtm-bff/internal/config"
	mockclients "github.com/DKhorkov/hmtm-bff/mocks/clients"
)

func TestNewS3FileStorageRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	logger := mocklogging.NewMockLogger(ctrl)

	testCases := []struct {
		name          string
		s3Config      appconfig.S3Config
		expectError   bool
		expectedError string
	}{
		{
			name: "success",
			s3Config: appconfig.S3Config{
				AccessKeyID:     "test-access-key",
				SecretAccessKey: "test-secret-key",
				Region:          "us-east-1",
				Bucket:          "test-bucket",
				ACL:             "public-read",
				Timeout:         5 * time.Second,
			},
			expectError:   false,
			expectedError: "",
		},
		// Тестирование ошибки конфигурации AWS сложнее из-за статического вызова config.LoadDefaultConfig,
		// поэтому этот случай опущен в примере, но может быть добавлен с использованием патчинга.
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo, err := NewS3FileStorageRepository(tc.s3Config, logger)
			if tc.expectError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedError)
				require.Nil(t, repo)
			} else {
				require.NoError(t, err)
				require.NotNil(t, repo)
				require.Equal(t, tc.s3Config, repo.s3config)
				require.NotNil(t, repo.client)
				require.Equal(t, logger, repo.logger)
			}
		})
	}
}

func TestS3FileStorageRepository_Upload(t *testing.T) {
	ctrl := gomock.NewController(t)
	logger := mocklogging.NewMockLogger(ctrl)
	s3Client := mockclients.NewMockS3Client(ctrl)

	s3Config := appconfig.S3Config{
		Bucket:  "test-bucket",
		Region:  "us-east-1",
		ACL:     "public-read",
		Timeout: 5 * time.Second,
	}

	repo := &S3FileStorageRepository{
		client:   s3Client,
		logger:   logger,
		s3config: s3Config,
	}

	testCases := []struct {
		name          string
		key           string
		file          []byte
		setupMocks    func(s3Client *mockclients.MockS3Client)
		expectedURL   string
		errorExpected bool
	}{
		{
			name: "success",
			key:  "test-key",
			file: []byte("test content"),
			setupMocks: func(s3Client *mockclients.MockS3Client) {
				s3Client.
					EXPECT().
					PutObject(
						gomock.Any(),
						&s3.PutObjectInput{
							Bucket: aws.String("test-bucket"),
							Key:    aws.String("test-key"),
							Body:   bytes.NewReader([]byte("test content")),
							ACL:    "public-read",
						},
						gomock.Any(),
					).
					Return(&s3.PutObjectOutput{}, nil).
					Times(1)
			},
			expectedURL:   "https://test-bucket.s3.us-east-1.amazonaws.com/test-key",
			errorExpected: false,
		},
		{
			name: "upload error",
			key:  "test-key",
			file: []byte("test content"),
			setupMocks: func(s3Client *mockclients.MockS3Client) {
				s3Client.
					EXPECT().
					PutObject(
						gomock.Any(),
						&s3.PutObjectInput{
							Bucket: aws.String("test-bucket"),
							Key:    aws.String("test-key"),
							Body:   bytes.NewReader([]byte("test content")),
							ACL:    "public-read",
						},
						gomock.Any(),
					).
					Return(nil, errors.New("upload failed")).
					Times(1)
			},
			expectedURL:   "",
			errorExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(s3Client)
			}

			url, err := repo.Upload(context.Background(), tc.key, tc.file)
			if tc.errorExpected {
				require.Error(t, err)
				require.Empty(t, url)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedURL, url)
			}
		})
	}
}

func TestS3FileStorageRepository_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	logger := mocklogging.NewMockLogger(ctrl)
	s3Client := mockclients.NewMockS3Client(ctrl)

	s3Config := appconfig.S3Config{
		Bucket:  "test-bucket",
		Timeout: 5 * time.Second,
	}

	repo := &S3FileStorageRepository{
		client:   s3Client,
		logger:   logger,
		s3config: s3Config,
	}

	testCases := []struct {
		name          string
		key           string
		setupMocks    func(s3Client *mockclients.MockS3Client)
		errorExpected bool
	}{
		{
			name: "success",
			key:  "test-key",
			setupMocks: func(s3Client *mockclients.MockS3Client) {
				s3Client.
					EXPECT().
					DeleteObject(
						gomock.Any(),
						&s3.DeleteObjectInput{
							Bucket: aws.String("test-bucket"),
							Key:    aws.String("test-key"),
						},
						gomock.Any(),
					).
					Return(&s3.DeleteObjectOutput{}, nil).
					Times(1)

				s3Client.
					EXPECT().
					HeadObject(
						gomock.Any(),
						&s3.HeadObjectInput{
							Bucket: aws.String("test-bucket"),
							Key:    aws.String("test-key"),
						},
						gomock.Any(),
					).
					Return(nil, &types.NotFound{}).
					Times(1)
			},
			errorExpected: false,
		},
		{
			name: "delete error",
			key:  "test-key",
			setupMocks: func(s3Client *mockclients.MockS3Client) {
				s3Client.
					EXPECT().
					DeleteObject(
						gomock.Any(),
						&s3.DeleteObjectInput{
							Bucket: aws.String("test-bucket"),
							Key:    aws.String("test-key"),
						},
						gomock.Any(),
					).
					Return(nil, errors.New("delete failed")).
					Times(1)
			},
			errorExpected: true,
		},
		{
			name: "wait error",
			key:  "test-key",
			setupMocks: func(s3Client *mockclients.MockS3Client) {
				s3Client.
					EXPECT().
					DeleteObject(
						gomock.Any(),
						&s3.DeleteObjectInput{
							Bucket: aws.String("test-bucket"),
							Key:    aws.String("test-key"),
						},
						gomock.Any(),
					).
					Return(&s3.DeleteObjectOutput{}, nil).
					Times(1)

				// Объект всё ещё существует после всех попыток
				s3Client.
					EXPECT().
					HeadObject(
						gomock.Any(),
						&s3.HeadObjectInput{
							Bucket: aws.String("test-bucket"),
							Key:    aws.String("test-key"),
						},
						gomock.Any(),
					).
					Return(&s3.HeadObjectOutput{}, nil).
					MinTimes(1) // Может быть вызван несколько раз до таймаута
			},
			errorExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(s3Client)
			}

			err := repo.Delete(context.Background(), tc.key)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestS3FileStorageRepository_DeleteMany(t *testing.T) {
	ctrl := gomock.NewController(t)
	logger := mocklogging.NewMockLogger(ctrl)
	s3Client := mockclients.NewMockS3Client(ctrl)

	s3Config := appconfig.S3Config{
		Bucket:  "test-bucket",
		Timeout: 5 * time.Second,
	}

	repo := &S3FileStorageRepository{
		client:   s3Client,
		logger:   logger,
		s3config: s3Config,
	}

	testCases := []struct {
		name           string
		keys           []string
		setupMocks     func(s3Client *mockclients.MockS3Client)
		expectedErrors int
	}{
		{
			name: "success",
			keys: []string{"key1", "key2"},
			setupMocks: func(s3Client *mockclients.MockS3Client) {
				s3Client.
					EXPECT().
					DeleteObjects(
						gomock.Any(),
						&s3.DeleteObjectsInput{
							Bucket: aws.String("test-bucket"),
							Delete: &types.Delete{
								Objects: []types.ObjectIdentifier{
									{Key: pointers.New("key1")},
									{Key: pointers.New("key2")},
								},
							},
						},
						gomock.Any(),
					).
					Return(&s3.DeleteObjectsOutput{
						Deleted: []types.DeletedObject{
							{Key: pointers.New("key1")},
							{Key: pointers.New("key2")},
						},
						Errors: []types.Error{},
					}, nil).
					Times(1)

				s3Client.
					EXPECT().
					HeadObject(
						gomock.Any(),
						&s3.HeadObjectInput{
							Bucket: aws.String("test-bucket"),
							Key:    pointers.New("key1"),
						},
						gomock.Any(),
					).
					Return(nil, &types.NotFound{}).
					Times(1)

				s3Client.
					EXPECT().
					HeadObject(
						gomock.Any(),
						&s3.HeadObjectInput{
							Bucket: aws.String("test-bucket"),
							Key:    pointers.New("key2"),
						},
						gomock.Any(),
					).
					Return(nil, &types.NotFound{}).
					Times(1)
			},
			expectedErrors: 0,
		},
		{
			name: "delete objects error",
			keys: []string{"key1"},
			setupMocks: func(s3Client *mockclients.MockS3Client) {
				s3Client.
					EXPECT().
					DeleteObjects(
						gomock.Any(),
						&s3.DeleteObjectsInput{
							Bucket: aws.String("test-bucket"),
							Delete: &types.Delete{
								Objects: []types.ObjectIdentifier{
									{Key: pointers.New("key1")},
								},
							},
						},
						gomock.Any(),
					).
					Return(nil, errors.New("delete failed")).
					Times(1)
			},
			expectedErrors: 1,
		},
		{
			name: "partial failure",
			keys: []string{"key1", "key2"},
			setupMocks: func(s3Client *mockclients.MockS3Client) {
				s3Client.
					EXPECT().
					DeleteObjects(
						gomock.Any(),
						&s3.DeleteObjectsInput{
							Bucket: aws.String("test-bucket"),
							Delete: &types.Delete{
								Objects: []types.ObjectIdentifier{
									{Key: pointers.New("key1")},
									{Key: pointers.New("key2")},
								},
							},
						},
						gomock.Any(),
					).
					Return(&s3.DeleteObjectsOutput{
						Deleted: []types.DeletedObject{
							{Key: pointers.New("key1")},
						},
						Errors: []types.Error{
							{Key: pointers.New("key2"), Message: aws.String("delete error")},
						},
					}, nil).
					Times(1)
			},
			expectedErrors: 1,
		},
		{
			name: "wait error",
			keys: []string{"key1"},
			setupMocks: func(s3Client *mockclients.MockS3Client) {
				s3Client.
					EXPECT().
					DeleteObjects(
						gomock.Any(),
						&s3.DeleteObjectsInput{
							Bucket: aws.String("test-bucket"),
							Delete: &types.Delete{
								Objects: []types.ObjectIdentifier{
									{Key: pointers.New("key1")},
								},
							},
						},
						gomock.Any(),
					).
					Return(&s3.DeleteObjectsOutput{
						Deleted: []types.DeletedObject{
							{Key: pointers.New("key1")},
						},
						Errors: []types.Error{},
					}, nil).
					Times(1)

				s3Client.
					EXPECT().
					HeadObject(
						gomock.Any(),
						&s3.HeadObjectInput{
							Bucket: aws.String("test-bucket"),
							Key:    pointers.New("key1"),
						},
						gomock.Any(),
					).
					Return(&s3.HeadObjectOutput{}, nil). // Объект всё ещё существует
					MinTimes(1)                          // Может быть вызван несколько раз до таймаута
			},
			expectedErrors: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(s3Client)
			}

			errs := repo.DeleteMany(context.Background(), tc.keys)
			require.Equal(t, tc.expectedErrors, len(errs))
		})
	}
}
