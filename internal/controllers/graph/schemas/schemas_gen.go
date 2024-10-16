// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package schemas

type LoginUserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Mutation struct {
}

type Query struct {
}

type RefreshTokensInput struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type RegisterUserInput struct {
	Credentials *LoginUserInput `json:"credentials"`
}
