package entities

import "time"

type RawAddToyDTO struct {
	AccessToken string   `json:"access_token"`
	CategoryID  uint32   `json:"category_id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float32  `json:"price"`
	Quantity    uint32   `json:"quantity"`
	TagIDs      []uint32 `json:"tag_ids"`
}

type AddToyDTO struct {
	UserID      uint64   `json:"user_id"`
	CategoryID  uint32   `json:"category_id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float32  `json:"price"`
	Quantity    uint32   `json:"quantity"`
	TagIDs      []uint32 `json:"tag_ids"`
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

type Toy struct {
	ID          uint64    `json:"id"`
	MasterID    uint64    `json:"master_id"`
	CategoryID  uint32    `json:"category_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float32   `json:"price"`
	Quantity    uint32    `json:"quantity"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Tags        []Tag     `json:"tags"`
}
