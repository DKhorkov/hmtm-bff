package repositories

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"github.com/DKhorkov/libs/pointers"

	appconfig "github.com/DKhorkov/hmtm-bff/internal/config"
)

func NewS3FileStorageRepository(
	s3config appconfig.S3Config,
	logger *slog.Logger,
) *S3FileStorageRepository {
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithCredentialsProvider(
			credentials.StaticCredentialsProvider{
				Value: aws.Credentials{
					AccessKeyID:     s3config.AccessKeyID,
					SecretAccessKey: s3config.SecretAccessKey,
				},
			},
		),
		config.WithRegion(s3config.Region),
	)

	if err != nil {
		log.Fatal(err)
	}

	// Create an Amazon S3 service client
	client := s3.NewFromConfig(cfg)

	return &S3FileStorageRepository{
		client:   client,
		logger:   logger,
		s3config: s3config,
	}
}

type S3FileStorageRepository struct {
	client   *s3.Client
	logger   *slog.Logger
	s3config appconfig.S3Config
}

func (repo *S3FileStorageRepository) Upload(ctx context.Context, key string, file []byte) (string, error) {
	_, err := repo.client.PutObject(
		ctx,
		&s3.PutObjectInput{
			Bucket: aws.String(repo.s3config.Bucket),
			Key:    aws.String(key),
			Body:   bytes.NewReader(file),
			ACL:    types.ObjectCannedACL(repo.s3config.ACL),
		},
	)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"https://%s.s3.%s.amazonaws.com/%s",
		repo.s3config.Bucket,
		repo.s3config.Region,
		key,
	), nil
}

func (repo *S3FileStorageRepository) Delete(ctx context.Context, key string) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(repo.s3config.Bucket),
		Key:    aws.String(key),
	}

	if _, err := repo.client.DeleteObject(ctx, input); err != nil {
		return err
	}

	if err := s3.NewObjectNotExistsWaiter(repo.client).Wait(
		ctx, &s3.HeadObjectInput{
			Bucket: aws.String(repo.s3config.Bucket),
			Key:    aws.String(key),
		},
		time.Minute,
	); err != nil {
		return err
	}

	return nil
}

func (repo *S3FileStorageRepository) DeleteMany(ctx context.Context, keys []string) []error {
	objectsToDelete := make([]types.ObjectIdentifier, 0, len(keys))
	for _, key := range keys {
		objectsToDelete = append(objectsToDelete, types.ObjectIdentifier{Key: pointers.Pointer(key)})
	}

	delOut, err := repo.client.DeleteObjects(
		ctx,
		&s3.DeleteObjectsInput{
			Bucket: aws.String(repo.s3config.Bucket),
			Delete: &types.Delete{
				Objects: objectsToDelete,
			},
		},
	)

	var out []error
	switch {
	case err != nil:
		out = append(out, err)
	case len(delOut.Errors) > 0:
		out = make([]error, 0, len(delOut.Errors))
		for _, err := range delOut.Errors {
			out = append(out, fmt.Errorf("%v-%v:%v-%v", err.VersionId, err.Code, err.Key, err.Message))
		}
	default:
		for _, delObj := range delOut.Deleted {
			if err = s3.NewObjectNotExistsWaiter(repo.client).Wait(
				ctx,
				&s3.HeadObjectInput{
					Bucket: aws.String(repo.s3config.Bucket),
					Key:    delObj.Key,
				},
				time.Minute,
			); err != nil {
				out = append(out, err)
			}
		}
	}

	return out
}
