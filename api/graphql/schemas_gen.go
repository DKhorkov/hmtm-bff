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

type VerifyUserEmailInput struct {
	VerifyEmailToken string `json:"verifyEmailToken"`
}
