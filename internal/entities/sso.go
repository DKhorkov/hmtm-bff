package entities

import "time"

type LoginUserDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokensDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RegisterUserDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	ID        uint64    `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
