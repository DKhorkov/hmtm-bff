package models

import "time"

type LoginUserDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokensDTO struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type RegisterUserDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	ID        uint64    `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
