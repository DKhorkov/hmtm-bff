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
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type RegisterUserDTO struct {
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
	Password    string `json:"password"`
}

type User struct {
	ID                uint64    `json:"id"`
	DisplayName       string    `json:"displayName"`
	Email             string    `json:"email"`
	EmailConfirmed    bool      `json:"emailConfirmed"`
	Password          string    `json:"password"`
	Phone             *string   `json:"phone,omitempty"`
	PhoneConfirmed    bool      `json:"phoneConfirmed"`
	Telegram          *string   `json:"telegram,omitempty"`
	TelegramConfirmed bool      `json:"telegramConfirmed"`
	Avatar            *string   `json:"avatar,omitempty"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

type RawUpdateUserProfileDTO struct {
	AccessToken string          `json:"accessToken"`
	DisplayName *string         `json:"displayName,omitempty"`
	Phone       *string         `json:"phone,omitempty"`
	Telegram    *string         `json:"telegram,omitempty"`
	Avatar      *graphql.Upload `json:"avatar,omitempty"`
}

type UpdateUserProfileDTO struct {
	AccessToken string  `json:"accessToken"`
	DisplayName *string `json:"displayName,omitempty"`
	Phone       *string `json:"phone,omitempty"`
	Telegram    *string `json:"telegram,omitempty"`
	Avatar      *string `json:"avatar,omitempty"`
}
