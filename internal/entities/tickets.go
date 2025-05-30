package entities

import (
	"time"

	"github.com/99designs/gqlgen/graphql"
)

type Respond struct {
	ID        uint64    `json:"id"`
	TicketID  uint64    `json:"tagliatelle"`
	MasterID  uint64    `json:"masterId"`
	Price     float32   `json:"price"`
	Comment   *string   `json:"comment,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type RespondToTicketDTO struct {
	TicketID uint64  `json:"ticketId"`
	UserID   uint64  `json:"userId"`
	Price    float32 `json:"price"`
	Comment  *string `json:"comment,omitempty"`
}

type RawRespondToTicketDTO struct {
	AccessToken string  `json:"accessToken"`
	TicketID    uint64  `json:"ticketId"`
	Price       float32 `json:"price"`
	Comment     *string `json:"comment,omitempty"`
}

type RawUpdateRespondDTO struct {
	AccessToken string   `json:"accessToken"`
	ID          uint64   `json:"id"`
	Price       *float32 `json:"price,omitempty"`
	Comment     *string  `json:"comment,omitempty"`
}

type UpdateRespondDTO struct {
	ID      uint64   `json:"id"`
	Price   *float32 `json:"price,omitempty"`
	Comment *string  `json:"comment,omitempty"`
}

type TicketAttachment struct {
	ID        uint64    `json:"id"`
	TicketID  uint64    `json:"ticketId"`
	Link      string    `json:"link"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type RawTicket struct {
	ID          uint64             `json:"id"`
	UserID      uint64             `json:"userId"`
	CategoryID  uint32             `json:"categoryId"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Price       *float32           `json:"price,omitempty"`
	Quantity    uint32             `json:"quantity"`
	CreatedAt   time.Time          `json:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt"`
	TagIDs      []uint32           `json:"tagIds"`
	Attachments []TicketAttachment `json:"attachments,omitempty"`
}

type Ticket struct {
	ID          uint64             `json:"id"`
	UserID      uint64             `json:"userId"`
	CategoryID  uint32             `json:"categoryId"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Price       *float32           `json:"price,omitempty"`
	Quantity    uint32             `json:"quantity"`
	CreatedAt   time.Time          `json:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt"`
	Tags        []Tag              `json:"tags,omitempty"`
	Attachments []TicketAttachment `json:"attachments,omitempty"`
}

type CreateTicketDTO struct {
	UserID      uint64   `json:"userId"`
	CategoryID  uint32   `json:"categoryId"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       *float32 `json:"price,omitempty"`
	Quantity    uint32   `json:"quantity"`
	TagIDs      []uint32 `json:"tagIds,omitempty"`
	Attachments []string `json:"attachments,omitempty"`
}

type RawCreateTicketDTO struct {
	AccessToken string            `json:"accessToken"`
	CategoryID  uint32            `json:"categoryId"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Price       *float32          `json:"price,omitempty"`
	Quantity    uint32            `json:"quantity"`
	Tags        []string          `json:"tags,omitempty"`
	Attachments []*graphql.Upload `json:"attachments,omitempty"`
}

type RawUpdateTicketDTO struct {
	AccessToken string            `json:"accessToken"`
	ID          uint64            `json:"id"`
	CategoryID  *uint32           `json:"categoryId,omitempty"`
	Name        *string           `json:"name,omitempty"`
	Description *string           `json:"description,omitempty"`
	Price       *float32          `json:"price,omitempty"`
	Quantity    *uint32           `json:"quantity,omitempty"`
	Tags        []string          `json:"tags,omitempty"`
	Attachments []*graphql.Upload `json:"attachments,omitempty"`
}

type UpdateTicketDTO struct {
	ID          uint64   `json:"id"`
	CategoryID  *uint32  `json:"categoryId,omitempty"`
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	Price       *float32 `json:"price,omitempty"`
	Quantity    *uint32  `json:"quantity,omitempty"`
	TagIDs      []uint32 `json:"tagIds,omitempty"`
	Attachments []string `json:"attachments,omitempty"`
}

type TicketsFilters struct {
	Search              *string  `json:"search,omitempty"`
	PriceCeil           *float32 `json:"priceCeil,omitempty"`     // max price
	PriceFloor          *float32 `json:"priceFloor,omitempty"`    // min price
	QuantityFloor       *uint32  `json:"quantityFloor,omitempty"` // min quantity
	CategoryIDs         []uint32 `json:"categoryIds,omitempty"`
	TagIDs              []uint32 `json:"tagIds,omitempty"`
	CreatedAtOrderByAsc *bool    `json:"createdAtOrderByAsc,omitempty"`
}
