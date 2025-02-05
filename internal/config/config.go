package config

import (
	"fmt"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/DKhorkov/libs/cookies"
	"github.com/DKhorkov/libs/loadenv"
	"github.com/DKhorkov/libs/logging"
	"github.com/DKhorkov/libs/tracing"
)

func New() Config {
	return Config{
		Environment: loadenv.GetEnv("ENVIRONMENT", "local"),
		Version:     loadenv.GetEnv("VERSION", "latest"),
		HTTP: HTTPConfig{
			Host: loadenv.GetEnv("HOST", "0.0.0.0"),
			Port: loadenv.GetEnvAsInt("PORT", 8080),
			ReadTimeout: time.Second * time.Duration(
				loadenv.GetEnvAsInt("HTTP_READ_TIMEOUT", 3),
			),
			ReadHeaderTimeout: time.Second * time.Duration(
				loadenv.GetEnvAsInt("HTTP_READ_HEADER_TIMEOUT", 1),
			),
			TimeoutHandlerTimeout: time.Second * time.Duration(
				loadenv.GetEnvAsInt("HTTP_TIMEOUT_HANDLER_TIMEOUT", 2),
			),
		},
		Clients: ClientsConfig{
			SSO: ClientConfig{
				Host:         loadenv.GetEnv("SSO_CLIENT_HOST", "0.0.0.0"),
				Port:         loadenv.GetEnvAsInt("SSO_CLIENT_PORT", 8070),
				RetriesCount: loadenv.GetEnvAsInt("SSO_RETRIES_COUNT", 3),
				RetryTimeout: time.Second * time.Duration(
					loadenv.GetEnvAsInt("SSO_RETRIES_TIMEOUT", 1),
				),
			},
			Toys: ClientConfig{
				Host:         loadenv.GetEnv("TOYS_CLIENT_HOST", "0.0.0.0"),
				Port:         loadenv.GetEnvAsInt("TOYS_CLIENT_PORT", 8060),
				RetriesCount: loadenv.GetEnvAsInt("TOYS_RETRIES_COUNT", 3),
				RetryTimeout: time.Second * time.Duration(
					loadenv.GetEnvAsInt("TOYS_RETRIES_TIMEOUT", 1),
				),
			},
			Tickets: ClientConfig{
				Host:         loadenv.GetEnv("TICKETS_CLIENT_HOST", "0.0.0.0"),
				Port:         loadenv.GetEnvAsInt("TICKETS_CLIENT_PORT", 8050),
				RetriesCount: loadenv.GetEnvAsInt("TICKETS_RETRIES_COUNT", 3),
				RetryTimeout: time.Second * time.Duration(
					loadenv.GetEnvAsInt("TICKETS_RETRIES_TIMEOUT", 1),
				),
			},
			Notifications: ClientConfig{
				Host:         loadenv.GetEnv("NOTIFICATIONS_CLIENT_HOST", "0.0.0.0"),
				Port:         loadenv.GetEnvAsInt("NOTIFICATIONS_CLIENT_HOST_CLIENT_PORT", 8040),
				RetriesCount: loadenv.GetEnvAsInt("NOTIFICATIONS_CLIENT_HOST_RETRIES_COUNT", 3),
				RetryTimeout: time.Second * time.Duration(
					loadenv.GetEnvAsInt("NOTIFICATIONS_CLIENT_HOST_RETRIES_TIMEOUT", 1),
				),
			},
		},
		Logging: logging.Config{
			Level:       logging.Levels.DEBUG,
			LogFilePath: fmt.Sprintf("logs/%s.log", time.Now().UTC().Format("02-01-2006")),
		},
		CORS: CORSConfig{
			AllowedOrigins:   loadenv.GetEnvAsSlice("CORS_ALLOWED_ORIGINS", []string{"*"}, ", "),
			AllowedMethods:   loadenv.GetEnvAsSlice("CORS_ALLOWED_METHODS", []string{"*"}, ", "),
			AllowedHeaders:   loadenv.GetEnvAsSlice("CORS_ALLOWED_HEADERS", []string{"*"}, ", "),
			AllowCredentials: loadenv.GetEnvAsBool("CORS_ALLOW_CREDENTIALS", true),
			MaxAge:           loadenv.GetEnvAsInt("CORS_MAX_AGE", 600),
		},
		Cookies: CookiesConfig{
			AccessToken: cookies.Config{
				Path:   loadenv.GetEnv("COOKIES_ACCESS_TOKEN_PATH", "/"),
				Domain: loadenv.GetEnv("COOKIES_ACCESS_TOKEN_DOMAIN", ""),
				MaxAge: loadenv.GetEnvAsInt("COOKIES_ACCESS_TOKEN_MAX_AGE", 0),
				Expires: time.Minute * time.Duration(
					loadenv.GetEnvAsInt("COOKIES_ACCESS_TOKEN_EXPIRES", 15),
				),
				Secure:   loadenv.GetEnvAsBool("COOKIES_ACCESS_TOKEN_SECURE", false),
				HTTPOnly: loadenv.GetEnvAsBool("COOKIES_ACCESS_TOKEN_HTTP_ONLY", false),
				SameSite: http.SameSite(
					loadenv.GetEnvAsInt("COOKIES_ACCESS_TOKEN_SAME_SITE", 1),
				),
			},
			RefreshToken: cookies.Config{
				Path:   loadenv.GetEnv("COOKIES_REFRESH_TOKEN_PATH", "/"),
				Domain: loadenv.GetEnv("COOKIES_REFRESH_TOKEN_DOMAIN", ""),
				MaxAge: loadenv.GetEnvAsInt("COOKIES_REFRESH_TOKEN_MAX_AGE", 0),
				Expires: time.Hour * time.Duration(
					loadenv.GetEnvAsInt("COOKIES_REFRESH_TOKEN_EXPIRES", 24*7),
				),
				Secure:   loadenv.GetEnvAsBool("COOKIES_REFRESH_TOKEN_SECURE", false),
				HTTPOnly: loadenv.GetEnvAsBool("COOKIES_REFRESH_TOKEN_HTTP_ONLY", false),
				SameSite: http.SameSite(
					loadenv.GetEnvAsInt("COOKIES_REFRESH_TOKEN_SAME_SITE", 1),
				),
			},
		},
		S3: S3Config{
			AccessKeyID:     loadenv.GetEnv("S3_ACCESS_KEY_ID", ""),
			SecretAccessKey: loadenv.GetEnv("S3_SECRET_ACCESS_KEY", ""),
			Region:          loadenv.GetEnv("S3_REGION", ""),
			Bucket:          loadenv.GetEnv("S3_BUCKET", ""),
			ACL:             loadenv.GetEnv("S3_ACL", "public-read"),
		},
		Validation: ValidationConfig{
			FileMaxSize: int64(loadenv.GetEnvAsInt("FILE_MAX_SIZE", 5*1024*1024)), // 5 Mb
			FileAllowedExtensions: loadenv.GetEnvAsSlice(
				"FILE_ALLOWED_EXTENSIONS",
				[]string{
					".png",
					".svg",
					".gif",
					".jpg",
					".jpeg",
					".jfif",
					".pjpeg",
					".pjp",
				},
				",",
			),
		},
		Tracing: TracingConfig{
			Server: tracing.Config{
				ServiceName:    loadenv.GetEnv("TRACING_SERVICE_NAME", "hmtm-bff"),
				ServiceVersion: loadenv.GetEnv("VERSION", "latest"),
				JaegerURL: fmt.Sprintf(
					"http://%s:%d/api/traces",
					loadenv.GetEnv("TRACING_JAEGER_HOST", "0.0.0.0"),
					loadenv.GetEnvAsInt("TRACING_API_TRACES_PORT", 14268),
				),
			},
			Spans: SpansConfig{
				Root: tracing.SpanConfig{
					Name: "Root",
					Opts: []trace.SpanStartOption{
						trace.WithAttributes(
							attribute.String("Environment", loadenv.GetEnv("ENVIRONMENT", "local")),
						),
					},
					Events: tracing.SpanEventsConfig{
						Start: tracing.SpanEventConfig{
							Name: "Calling handler",
							Opts: []trace.EventOption{
								trace.WithAttributes(
									attribute.String("Environment", loadenv.GetEnv("ENVIRONMENT", "local")),
								),
							},
						},
						End: tracing.SpanEventConfig{
							Name: "Received response from handler",
							Opts: []trace.EventOption{
								trace.WithAttributes(
									attribute.String("Environment", loadenv.GetEnv("ENVIRONMENT", "local")),
								),
							},
						},
					},
				},
				Clients: SpanClients{
					SSO: tracing.SpanConfig{
						Opts: []trace.SpanStartOption{
							trace.WithAttributes(
								attribute.String("Environment", loadenv.GetEnv("ENVIRONMENT", "local")),
							),
						},
						Events: tracing.SpanEventsConfig{
							Start: tracing.SpanEventConfig{
								Name: "Calling gRPC SSO client",
								Opts: []trace.EventOption{
									trace.WithAttributes(
										attribute.String("Environment", loadenv.GetEnv("ENVIRONMENT", "local")),
									),
								},
							},
							End: tracing.SpanEventConfig{
								Name: "Received response from gRPC SSO client",
								Opts: []trace.EventOption{
									trace.WithAttributes(
										attribute.String("Environment", loadenv.GetEnv("ENVIRONMENT", "local")),
									),
								},
							},
						},
					},
					Toys: tracing.SpanConfig{
						Opts: []trace.SpanStartOption{
							trace.WithAttributes(
								attribute.String("Environment", loadenv.GetEnv("ENVIRONMENT", "local")),
							),
						},
						Events: tracing.SpanEventsConfig{
							Start: tracing.SpanEventConfig{
								Name: "Calling gRPC Toys client",
								Opts: []trace.EventOption{
									trace.WithAttributes(
										attribute.String("Environment", loadenv.GetEnv("ENVIRONMENT", "local")),
									),
								},
							},
							End: tracing.SpanEventConfig{
								Name: "Received response from gRPC Toys client",
								Opts: []trace.EventOption{
									trace.WithAttributes(
										attribute.String("Environment", loadenv.GetEnv("ENVIRONMENT", "local")),
									),
								},
							},
						},
					},
					Tickets: tracing.SpanConfig{
						Opts: []trace.SpanStartOption{
							trace.WithAttributes(
								attribute.String("Environment", loadenv.GetEnv("ENVIRONMENT", "local")),
							),
						},
						Events: tracing.SpanEventsConfig{
							Start: tracing.SpanEventConfig{
								Name: "Calling gRPC Tickets client",
								Opts: []trace.EventOption{
									trace.WithAttributes(
										attribute.String("Environment", loadenv.GetEnv("ENVIRONMENT", "local")),
									),
								},
							},
							End: tracing.SpanEventConfig{
								Name: "Received response from gRPC Tickets client",
								Opts: []trace.EventOption{
									trace.WithAttributes(
										attribute.String("Environment", loadenv.GetEnv("ENVIRONMENT", "local")),
									),
								},
							},
						},
					},
					Notifications: tracing.SpanConfig{
						Opts: []trace.SpanStartOption{
							trace.WithAttributes(
								attribute.String("Environment", loadenv.GetEnv("ENVIRONMENT", "local")),
							),
						},
						Events: tracing.SpanEventsConfig{
							Start: tracing.SpanEventConfig{
								Name: "Calling gRPC Notifications client",
								Opts: []trace.EventOption{
									trace.WithAttributes(
										attribute.String("Environment", loadenv.GetEnv("ENVIRONMENT", "local")),
									),
								},
							},
							End: tracing.SpanEventConfig{
								Name: "Received response from gRPC Notifications client",
								Opts: []trace.EventOption{
									trace.WithAttributes(
										attribute.String("Environment", loadenv.GetEnv("ENVIRONMENT", "local")),
									),
								},
							},
						},
					},
				},
			},
		},
	}
}

type HTTPConfig struct {
	Host                  string
	Port                  int
	ReadHeaderTimeout     time.Duration
	ReadTimeout           time.Duration
	TimeoutHandlerTimeout time.Duration
}

type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	MaxAge           int
	AllowCredentials bool
}

type ClientConfig struct {
	Host         string
	Port         int
	RetryTimeout time.Duration
	RetriesCount int
}

type ClientsConfig struct {
	SSO           ClientConfig
	Toys          ClientConfig
	Tickets       ClientConfig
	Notifications ClientConfig
}

type CookiesConfig struct {
	AccessToken  cookies.Config
	RefreshToken cookies.Config
}

type S3Config struct {
	Region          string
	Bucket          string
	AccessKeyID     string
	SecretAccessKey string
	ACL             string
}

type ValidationConfig struct {
	FileMaxSize           int64
	FileAllowedExtensions []string
}

type TracingConfig struct {
	Server tracing.Config
	Spans  SpansConfig
}

type SpansConfig struct {
	Root    tracing.SpanConfig
	Clients SpanClients
}

type SpanClients struct {
	SSO           tracing.SpanConfig
	Toys          tracing.SpanConfig
	Tickets       tracing.SpanConfig
	Notifications tracing.SpanConfig
}

type Config struct {
	HTTP        HTTPConfig
	CORS        CORSConfig
	Clients     ClientsConfig
	Logging     logging.Config
	Cookies     CookiesConfig
	S3          S3Config
	Validation  ValidationConfig
	Tracing     TracingConfig
	Environment string
	Version     string
}
