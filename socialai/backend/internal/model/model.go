package model

import "time"

type User struct {
	Username     string    `json:"username"`
	PasswordHash string    `json:"password_hash,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

type Post struct {
	ID        string    `json:"id"`
	User      string    `json:"user"`
	Message   string    `json:"message"`
	URL       string    `json:"url"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}

type SearchFilter struct {
	User     string
	Keywords string
	Limit    int
}
