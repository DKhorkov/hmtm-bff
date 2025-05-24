package entities

import (
	"time"

	"github.com/99designs/gqlgen/graphql"
)

type RawAddToyDTO struct {
	AccessToken string            `json:"accessToken"`
	CategoryID  uint32            `json:"categoryId"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Price       float32           `json:"price"`
	Quantity    uint32            `json:"quantity"`
	Tags        []string          `json:"tags,omitempty"`
	Attachments []*graphql.Upload `json:"attachments,omitempty"`
}

type AddToyDTO struct {
	UserID      uint64   `json:"userId"`
	CategoryID  uint32   `json:"categoryId"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float32  `json:"price"`
	Quantity    uint32   `json:"quantity"`
	TagIDs      []uint32 `json:"tagIds,omitempty"`
	Attachments []string `json:"attachments,omitempty"`
}

type Category struct {
	ID   uint32 `json:"id"`
	Name string `json:"name"`
}

type Master struct {
	ID        uint64    `json:"id"`
	UserID    uint64    `json:"userId"`
	Info      *string   `json:"info,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type RawRegisterMasterDTO struct {
	AccessToken string  `json:"accessToken"`
	Info        *string `json:"info,omitempty"`
}

type RegisterMasterDTO struct {
	UserID uint64  `json:"userId"`
	Info   *string `json:"info,omitempty"`
}

type RawUpdateMasterDTO struct {
	AccessToken string  `json:"accessToken"`
	ID          uint64  `json:"id"`
	Info        *string `json:"info,omitempty"`
}

type UpdateMasterDTO struct {
	ID   uint64  `json:"id"`
	Info *string `json:"info,omitempty"`
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
	ToyID     uint64    `json:"toyId"`
	Link      string    `json:"link"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Toy struct {
	ID          uint64          `json:"id"`
	MasterID    uint64          `json:"masterId"`
	CategoryID  uint32          `json:"categoryId"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Price       float32         `json:"price"`
	Quantity    uint32          `json:"quantity"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
	Tags        []Tag           `json:"tags,omitempty"`
	Attachments []ToyAttachment `json:"attachments,omitempty"`
}

type RawUpdateToyDTO struct {
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

type UpdateToyDTO struct {
	ID          uint64   `json:"id"`
	CategoryID  *uint32  `json:"categoryId,omitempty"`
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	Price       *float32 `json:"price,omitempty"`
	Quantity    *uint32  `json:"quantity,omitempty"`
	TagIDs      []uint32 `json:"tagIds,omitempty"`
	Attachments []string `json:"attachments,omitempty"`
}

type ToysFilters struct {
	Search              *string  `json:"search,omitempty"`
	PriceCeil           *float32 `json:"priceCeil,omitempty"`     // max price
	PriceFloor          *float32 `json:"priceFloor,omitempty"`    // min price
	QuantityFloor       *uint32  `json:"quantityFloor,omitempty"` // min quantity
	CategoryIDs         []uint32 `json:"categoryIds,omitempty"`
	TagIDs              []uint32 `json:"tagIds,omitempty"`
	CreatedAtOrderByAsc *bool    `json:"createdAtOrderByAsc,omitempty"`
}
