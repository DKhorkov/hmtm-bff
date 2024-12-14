package config

import (
	"fmt"
	"net/http"
	"time"

	"github.com/DKhorkov/libs/cookies"
	"github.com/DKhorkov/libs/loadenv"
	"github.com/DKhorkov/libs/logging"
)

func New() Config {
	return Config{
		HTTP: HTTPConfig{
			Host: loadenv.GetEnv("HOST", "0.0.0.0"),
			Port: loadenv.GetEnvAsInt("PORT", 8080),
			ReadHeaderTimeout: time.Second * time.Duration(
				loadenv.GetEnvAsInt("HTTP_READ_HEADER_TIMEOUT", 1),
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
		},
		Logging: logging.Config{
			Level:       logging.Levels.DEBUG,
			LogFilePath: fmt.Sprintf("logs/%s.log", time.Now().Format("02-01-2006")),
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

type ClientConfig struct {
	Host         string
	Port         int
	RetryTimeout time.Duration
	RetriesCount int
}

type ClientsConfig struct {
	SSO  ClientConfig
	Toys ClientConfig
}

type CookiesConfig struct {
	AccessToken  cookies.Config
	RefreshToken cookies.Config
}

type Config struct {
	HTTP    HTTPConfig
	CORS    CORSConfig
	Clients ClientsConfig
	Logging logging.Config
	Cookies CookiesConfig
}
