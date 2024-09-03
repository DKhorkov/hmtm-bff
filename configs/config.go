package configs

import "hmtmbff/pkg/loadenv"

func GetConfig() *Config {
	return &Config{
		Graphql: GraphqlConfigs{
			Port: loadenv.GetEnvAsInt("GRAPHQL_PORT", 8080),
		},
	}
}

type GraphqlConfigs struct {
	Port int
}

type Config struct {
	Graphql GraphqlConfigs
}
