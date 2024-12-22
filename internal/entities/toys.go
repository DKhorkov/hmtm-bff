package entities

import "time"

type AddToyDTO struct {
	AccessToken string   `json:"accessToken"`
	CategoryID  uint32   `json:"categoryID"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float32  `json:"price"`
	Quantity    uint32   `json:"quantity"`
	TagsIDs     []uint32 `json:"tags"`
}

type Category struct {
	ID   uint32 `json:"id"`
	Name string `json:"name"`
}

type Master struct {
	ID        uint64    `json:"id"`
	UserID    uint64    `json:"userID"`
	Info      string    `json:"info"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type RegisterMasterDTO struct {
	AccessToken string `json:"accessToken"`
	Info        string `json:"info"`
}

type Tag struct {
	ID   uint32 `json:"id"`
	Name string `json:"name"`
}

type Toy struct {
	ID          uint64    `json:"id"`
	MasterID    uint64    `json:"masterID"`
	CategoryID  uint32    `json:"categoryID"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float32   `json:"price"`
	Quantity    uint32    `json:"quantity"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Tags        []Tag     `json:"tags"`
}
