package repositories

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	appconfig "github.com/DKhorkov/hmtm-bff/internal/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
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
