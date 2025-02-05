package entities

import "time"

type Email struct {
	ID      uint64    `json:"id"`
	UserID  uint64    `json:"user_id"`
	Email   string    `json:"email"`
	Content string    `json:"content"`
	SentAt  time.Time `json:"sent_at"`
}
