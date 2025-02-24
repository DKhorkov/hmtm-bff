// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package graphqlapi

import (
	"github.com/99designs/gqlgen/graphql"
)

type AddToyInput struct {
	CategoryID  string            `json:"categoryId"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Price       float64           `json:"price"`
	Quantity    int               `json:"quantity"`
	Tags        []string          `json:"tags,omitempty"`
	Attachments []*graphql.Upload `json:"attachments,omitempty"`
}

type ChangePasswordInput struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

type CreateTicketInput struct {
	CategoryID  string            `json:"categoryId"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Price       float64           `json:"price"`
	Quantity    int               `json:"quantity"`
	Tags        []string          `json:"tags,omitempty"`
	Attachments []*graphql.Upload `json:"attachments,omitempty"`
}

type LoginUserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Mutation struct {
}

type Query struct {
}

type RegisterMasterInput struct {
	Info string `json:"info"`
}

type RegisterUserInput struct {
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
	Password    string `json:"password"`
}

type RespondToTicketInput struct {
	TicketID string `json:"ticketId"`
}

type SendVerifyEmailMessageInput struct {
	Email string `json:"email"`
}

type UpdateUserProfileInput struct {
	DisplayName *string         `json:"displayName,omitempty"`
	Phone       *string         `json:"phone,omitempty"`
	Telegram    *string         `json:"telegram,omitempty"`
	Avatar      *graphql.Upload `json:"avatar,omitempty"`
}

type VerifyUserEmailInput struct {
	VerifyEmailToken string `json:"verifyEmailToken"`
}
