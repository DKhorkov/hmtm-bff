package loadenv

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

// GetEnv is a helper function to read an environment or return a default value.
func GetEnv(key string, defaultVal string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}

	return defaultVal
}

// GetEnvAsInt is a helper function to read an environment variable into integer or return a default value.
func GetEnvAsInt(name string, defaultVal int) int {
	valueStr := GetEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}

// GetEnvAsBool is a helper to read an environment variable into a bool or return default value.
func GetEnvAsBool(name string, defaultVal bool) bool {
	valStr := GetEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultVal
}

// GetEnvAsSlice is a helper to read an environment variable into a string slice or return default value.
func GetEnvAsSlice(name string, defaultVal []string, sep string) []string {
	valStr := GetEnv(name, "")

	if valStr == "" {
		return defaultVal
	}

	val := strings.Split(valStr, sep)
	return val
}
