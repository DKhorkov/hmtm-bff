package interfaces

import (
	"context"

	"github.com/DKhorkov/hmtm-notifications/api/protobuf/generated/go/notifications"
	"github.com/DKhorkov/hmtm-sso/api/protobuf/generated/go/sso"
	"github.com/DKhorkov/hmtm-tickets/api/protobuf/generated/go/tickets"
	"github.com/DKhorkov/hmtm-toys/api/protobuf/generated/go/toys"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

//go:generate mockgen -source=clients.go -destination=../../mocks/clients/sso_client.go -package=mockclients -exclude_interfaces=ToysClient,TicketsClient,NotificationsClient,S3Client
type SsoClient interface {
	sso.AuthServiceClient
	sso.UsersServiceClient
}

//go:generate mockgen -source=clients.go -destination=../../mocks/clients/toys_client.go -package=mockclients -exclude_interfaces=SsoClient,TicketsClient,NotificationsClient,S3Client
type ToysClient interface {
	toys.CategoriesServiceClient
	toys.ToysServiceClient
	toys.TagsServiceClient
	toys.MastersServiceClient
}

//go:generate mockgen -source=clients.go -destination=../../mocks/clients/tickets_client.go -package=mockclients -exclude_interfaces=ToysClient,SsoClient,NotificationsClient,S3Client
type TicketsClient interface {
	tickets.TicketsServiceClient
	tickets.RespondsServiceClient
}

//go:generate mockgen -source=clients.go -destination=../../mocks/clients/notifications_client.go -package=mockclients -exclude_interfaces=ToysClient,TicketsClient,SsoClient,S3Client
type NotificationsClient interface {
	notifications.EmailsServiceClient
}

//go:generate mockgen -source=clients.go -destination=../../mocks/clients/s3_client.go -package=mockclients -exclude_interfaces=ToysClient,TicketsClient,SsoClient,NotificationsClient
type S3Client interface {
	PutObject(
		ctx context.Context,
		params *s3.PutObjectInput,
		optFns ...func(*s3.Options),
	) (*s3.PutObjectOutput, error)

	DeleteObject(
		ctx context.Context,
		params *s3.DeleteObjectInput,
		optFns ...func(*s3.Options),
	) (*s3.DeleteObjectOutput, error)

	DeleteObjects(
		ctx context.Context,
		params *s3.DeleteObjectsInput,
		optFns ...func(*s3.Options),
	) (*s3.DeleteObjectsOutput, error)

	HeadObject(
		ctx context.Context,
		params *s3.HeadObjectInput,
		optFns ...func(*s3.Options),
	) (*s3.HeadObjectOutput, error)
}
