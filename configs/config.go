package configs

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

/*
init is invoked before main()

https://pkg.go.dev/github.com/joho/godotenv#section-readme
https://habr.com/ru/articles/446468/
*/
func init() {
	// loads values from .env into the system.
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func GetConfig() *Config {
	return &Config{
		Graphql: GraphqlConfigs{
			Port: getEnvAsInt("GRAPHQL_PORT", 8080),
		},
	}
}

type GraphqlConfigs struct {
	Port int `env:"PORT"`
}

type Config struct {
	Graphql GraphqlConfigs
}

// Simple helper function to read an environment or return a default value.
func getEnv(key string, defaultVal string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}

	return defaultVal
}

// Simple helper function to read an environment variable into integer or return a default value.
func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}

// Helper to read an environment variable into a bool or return default value.
func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultVal
}

// Helper to read an environment variable into a string slice or return default value.
func getEnvAsSlice(name string, defaultVal []string, sep string) []string {
	valStr := getEnv(name, "")

	if valStr == "" {
		return defaultVal
	}

	val := strings.Split(valStr, sep)
	return val
}
