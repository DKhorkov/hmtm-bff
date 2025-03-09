package entities

import (
	"time"

	"github.com/99designs/gqlgen/graphql"
)

type RawAddToyDTO struct {
	AccessToken string            `json:"access_token"`
	CategoryID  uint32            `json:"category_id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Price       float32           `json:"price"`
	Quantity    uint32            `json:"quantity"`
	Tags        []string          `json:"tags,omitempty"`
	Attachments []*graphql.Upload `json:"attachments,omitempty"`
}

type AddToyDTO struct {
	UserID      uint64   `json:"user_id"`
	CategoryID  uint32   `json:"category_id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float32  `json:"price"`
	Quantity    uint32   `json:"quantity"`
	TagIDs      []uint32 `json:"tag_ids,omitempty"`
	Attachments []string `json:"attachments,omitempty"`
}

type Category struct {
	ID   uint32 `json:"id"`
	Name string `json:"name"`
}

type Master struct {
	ID        uint64    `json:"id"`
	UserID    uint64    `json:"user_id"`
	Info      string    `json:"info"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RawRegisterMasterDTO struct {
	AccessToken string `json:"access_token"`
	Info        string `json:"info"`
}

type RegisterMasterDTO struct {
	UserID uint64 `json:"user_id"`
	Info   string `json:"info"`
}

type Tag struct {
	ID   uint32 `json:"id"`
	Name string `json:"name"`
}

type CreateTagDTO struct {
	Name string `json:"name"`
}

type ToyAttachment struct {
	ID        uint64    `json:"id"`
	ToyID     uint64    `json:"toy_id"`
	Link      string    `json:"link"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Toy struct {
	ID          uint64          `json:"id"`
	MasterID    uint64          `json:"master_id"`
	CategoryID  uint32          `json:"category_id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Price       float32         `json:"price"`
	Quantity    uint32          `json:"quantity"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	Tags        []Tag           `json:"tags,omitempty"`
	Attachments []ToyAttachment `json:"attachments,omitempty"`
}

type RawUpdateToyDTO struct {
	AccessToken string            `json:"access_token"`
	ID          uint64            `json:"id"`
	CategoryID  uint32            `json:"category_id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Price       float32           `json:"price"`
	Quantity    uint32            `json:"quantity"`
	Tags        []string          `json:"tags,omitempty"`
	Attachments []*graphql.Upload `json:"attachments,omitempty"`
}

type UpdateToyDTO struct {
	ID          uint64   `json:"id"`
	CategoryID  uint32   `json:"category_id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float32  `json:"price"`
	Quantity    uint32   `json:"quantity"`
	TagIDs      []uint32 `json:"tag_ids,omitempty"`
	Attachments []string `json:"attachments,omitempty"`
}
