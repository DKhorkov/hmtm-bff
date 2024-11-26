package config

import (
	"fmt"
	"time"

	"github.com/DKhorkov/libs/loadenv"
	"github.com/DKhorkov/libs/logging"
)

func New() *Config {
	return &Config{
		HTTP: HTTPConfig{
			Host: loadenv.GetEnv("HOST", "0.0.0.0"),
			Port: loadenv.GetEnvAsInt("PORT", 8080),
			ReadHeaderTimeout: time.Second * time.Duration(
				loadenv.GetEnvAsInt("HTTP_READ_HEADER_TIMEOUT", 1),
			),
		},
		Clients: ClientsConfig{
			SSO: Client{
				Host:         loadenv.GetEnv("SSO_CLIENT_HOST", "0.0.0.0"),
				Port:         loadenv.GetEnvAsInt("SSO_CLIENT_PORT", 8070),
				RetriesCount: loadenv.GetEnvAsInt("SSO_RETRIES_COUNT", 3),
				RetryTimeout: time.Second * time.Duration(
					loadenv.GetEnvAsInt("SSO_RETRIES_TIMEOUT", 1),
				),
			},
		},
		Logging: logging.Config{
			Level:       logging.LogLevels.DEBUG,
			LogFilePath: fmt.Sprintf("logs/%s.log", time.Now().Format("02-01-2006")),
		},
		CORS: CORSConfig{
			AllowedOrigins:   loadenv.GetEnvAsSlice("CORS_ALLOWED_ORIGINS", []string{"*"}, ", "),
			AllowedMethods:   loadenv.GetEnvAsSlice("CORS_ALLOWED_METHODS", []string{"*"}, ", "),
			AllowedHeaders:   loadenv.GetEnvAsSlice("CORS_ALLOWED_HEADERS", []string{"*"}, ", "),
			AllowCredentials: loadenv.GetEnvAsBool("CORS_ALLOW_CREDENTIALS", true),
			MaxAge:           loadenv.GetEnvAsInt("CORS_MAX_AGE", 600),
		},
	}
}

type HTTPConfig struct {
	Host              string
	Port              int
	ReadHeaderTimeout time.Duration
}

type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	MaxAge           int
	AllowCredentials bool
}

type Client struct {
	Host         string
	Port         int
	RetryTimeout time.Duration
	RetriesCount int
}

type ClientsConfig struct {
	SSO Client
}

type Config struct {
	HTTP    HTTPConfig
	CORS    CORSConfig
	Clients ClientsConfig
	Logging logging.Config
}
