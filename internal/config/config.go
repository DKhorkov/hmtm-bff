package config

import (
	"time"

	"github.com/DKhorkov/hmtm-bff/pkg/loadenv"
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
	}
}

type HTTPConfig struct {
	Host              string
	Port              int
	ReadHeaderTimeout time.Duration
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
	Clients ClientsConfig
}
