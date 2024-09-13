package config

import "github.com/DKhorkov/hmtm-bff/pkg/loadenv"

func GetConfig() *Config {
	return &Config{
		HTTP: HTTPConfigs{
			Host:              loadenv.GetEnv("GRAPHQL_HOST", "0.0.0.0"),
			Port:              loadenv.GetEnvAsInt("GRAPHQL_PORT", 8080),
			ReadHeaderTimeout: loadenv.GetEnvAsInt("HTTP_READ_HEADER_TIMEOUT", 1),
		},
	}
}

type HTTPConfigs struct {
	Host              string
	Port              int
	ReadHeaderTimeout int // in seconds
}

type Config struct {
	HTTP HTTPConfigs
}
