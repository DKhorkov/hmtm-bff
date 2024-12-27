package ssogrpcclient

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/DKhorkov/hmtm-sso/api/protobuf/generated/go/sso"
	customgrpc "github.com/DKhorkov/libs/grpc"
	"github.com/DKhorkov/libs/logging"
	grpclogging "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	sso.AuthServiceClient
	sso.UsersServiceClient
}

func New(
	host string,
	port int,
	retriesCount int,
	retriesTimeout time.Duration,
	logger *slog.Logger,
) (*Client, error) {
	// Options for interceptors (перехватчики / middlewares) for retries purposes:
	retryOptions := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
		grpcretry.WithMax(uint(retriesCount)),
		grpcretry.WithPerRetryTimeout(retriesTimeout),
	}

	// Options for interceptors for logging purposes:
	logOptions := []grpclogging.Option{
		grpclogging.WithLogOnEvents(
			grpclogging.PayloadReceived,
			grpclogging.PayloadSent,
		),
	}

	// Create connection with SSO gRPC-server for client:
	clientConnection, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", host, port),
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
		grpc.WithChainUnaryInterceptor( // Middlewares. Using chain not to overwrite interceptors.
			grpclogging.UnaryClientInterceptor(
				customgrpc.UnaryClientLoggingInterceptor(logger),
				logOptions...,
			),
			grpcretry.UnaryClientInterceptor(retryOptions...),
		),
	)

	if err != nil {
		logging.LogError(
			logger,
			"Failed to create SSO gRPC client",
			err,
		)

		return nil, err
	}

	return &Client{
		sso.NewAuthServiceClient(clientConnection),
		sso.NewUsersServiceClient(clientConnection),
	}, nil
}
