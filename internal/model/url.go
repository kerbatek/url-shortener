package model

import "time"

type URL struct {
	ID          string    `json:"id" db:"id"`
	Code        string    `json:"code" db:"code"`
	OriginalURL string    `json:"original_url" db:"original_url"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
