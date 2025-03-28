package entities

import (
	"time"

	"github.com/99designs/gqlgen/graphql"
)

type LoginUserDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokensDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RegisterUserDTO struct {
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
}

type User struct {
	ID                uint64    `json:"id"`
	DisplayName       string    `json:"display_name"`
	Email             string    `json:"email"`
	EmailConfirmed    bool      `json:"email_confirmed"`
	Password          string    `json:"password"`
	Phone             *string   `json:"phone,omitempty"`
	PhoneConfirmed    bool      `json:"phone_confirmed"`
	Telegram          *string   `json:"telegram,omitempty"`
	TelegramConfirmed bool      `json:"telegram_confirmed"`
	Avatar            *string   `json:"avatar,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type RawUpdateUserProfileDTO struct {
	AccessToken string          `json:"access_token"`
	DisplayName *string         `json:"display_name,omitempty"`
	Phone       *string         `json:"phone,omitempty"`
	Telegram    *string         `json:"telegram,omitempty"`
	Avatar      *graphql.Upload `json:"avatar,omitempty"`
}

type UpdateUserProfileDTO struct {
	AccessToken string  `json:"access_token"`
	DisplayName *string `json:"display_name,omitempty"`
	Phone       *string `json:"phone,omitempty"`
	Telegram    *string `json:"telegram,omitempty"`
	Avatar      *string `json:"avatar,omitempty"`
}
