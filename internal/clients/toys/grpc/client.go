package toysgrpcclient

import (
	"fmt"
	"time"

	"github.com/DKhorkov/hmtm-toys/api/protobuf/generated/go/toys"
	"github.com/DKhorkov/libs/logging"
	"github.com/DKhorkov/libs/tracing"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"

	customgrpc "github.com/DKhorkov/libs/grpc/interceptors"
	grpclogging "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
)

type Client struct {
	toys.ToysServiceClient
	toys.TagsServiceClient
	toys.MastersServiceClient
	toys.CategoriesServiceClient
}

func New(
	host string,
	port int,
	retriesCount int,
	retriesTimeout time.Duration,
	logger logging.Logger,
	traceProvider tracing.Provider,
	spanConfig tracing.SpanConfig,
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
			customgrpc.UnaryClientTracingInterceptor(traceProvider, spanConfig),
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
			"Failed to create Toys gRPC client",
			err,
		)

		return nil, err
	}

	return &Client{
		ToysServiceClient:       toys.NewToysServiceClient(clientConnection),
		TagsServiceClient:       toys.NewTagsServiceClient(clientConnection),
		MastersServiceClient:    toys.NewMastersServiceClient(clientConnection),
		CategoriesServiceClient: toys.NewCategoriesServiceClient(clientConnection),
	}, nil
}
