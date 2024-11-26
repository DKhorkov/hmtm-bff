package ssogrpcclient

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/DKhorkov/hmtm-sso/protobuf/generated/go/sso"
	"github.com/DKhorkov/libs/logging"
	grpclogging "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	Auth  sso.AuthServiceClient
	Users sso.UsersServiceClient
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
			grpclogging.UnaryClientInterceptor(InterceptorLogger(logger), logOptions...),
			grpcretry.UnaryClientInterceptor(retryOptions...),
		),
	)

	if err != nil {
		logger.Error(
			"Failed to create SSO gRPC client",
			"Traceback",
			logging.GetLogTraceback(),
			"Error",
			err,
		)

		return nil, err
	}

	return &Client{
		Auth:  sso.NewAuthServiceClient(clientConnection),
		Users: sso.NewUsersServiceClient(clientConnection),
	}, nil
}

// InterceptorLogger adapts slog logger to interceptor logger.
// This code is simple enough to be copied and not imported.
func InterceptorLogger(logger *slog.Logger) grpclogging.Logger {
	return grpclogging.LoggerFunc(
		func(
			ctx context.Context,
			logLevel grpclogging.Level,
			msg string,
			fields ...any,
		) {
			logger.Log(
				ctx,
				slog.Level(logLevel),
				msg,
				fields...,
			)
		},
	)
}
