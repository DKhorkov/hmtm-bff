package entities

import "time"

type Respond struct {
	ID        uint64    `json:"id"`
	TicketID  uint64    `json:"ticket_id"`
	MasterID  uint64    `json:"master_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RespondToTicketDTO struct {
	TicketID uint64 `json:"ticket_id"`
	UserID   uint64 `json:"user_id"`
}

type RawRespondToTicketDTO struct {
	TicketID    uint64 `json:"ticket_id"`
	AccessToken string `json:"access_token"`
}

type RawTicket struct {
	ID          uint64    `json:"id"`
	UserID      uint64    `json:"user_id"`
	CategoryID  uint32    `json:"category_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float32   `json:"price"`
	Quantity    uint32    `json:"quantity"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	TagIDs      []uint32  `json:"tag_ids"`
}

type Ticket struct {
	ID          uint64    `json:"id"`
	UserID      uint64    `json:"user_id"`
	CategoryID  uint32    `json:"category_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float32   `json:"price"`
	Quantity    uint32    `json:"quantity"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Tags        []Tag     `json:"tags"`
}

type CreateTicketDTO struct {
	UserID      uint64   `json:"user_id"`
	CategoryID  uint32   `json:"category_id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float32  `json:"price"`
	Quantity    uint32   `json:"quantity"`
	TagIDs      []uint32 `json:"tag_ids"`
}

type RawCreateTicketDTO struct {
	AccessToken string   `json:"access_token"`
	CategoryID  uint32   `json:"category_id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float32  `json:"price"`
	Quantity    uint32   `json:"quantity"`
	TagIDs      []uint32 `json:"tag_ids"`
}
