package configs

import "github.com/DKhorkov/hmtm-bff/pkg/loadenv"

func GetConfig() *Config {
	return &Config{
		Graphql: GraphqlConfigs{
			Port: loadenv.GetEnvAsInt("GRAPHQL_PORT", 8080),
		},
		HTTP: HTTPConfigs{
			ReadHeaderTimeout: loadenv.GetEnvAsInt("HTTP_READ_HEADER_TIMEOUT", 1),
		},
	}
}

type GraphqlConfigs struct {
	Port int
}

type HTTPConfigs struct {
	ReadHeaderTimeout int // in seconds
}

type Config struct {
	Graphql GraphqlConfigs
	HTTP    HTTPConfigs
}
