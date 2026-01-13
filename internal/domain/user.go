package domain

import "time"

type User struct {
	ID        uint64     `json:"id"`
	Email     string     `json:"email"`
	Username  string     `json:"username"`
	Password  string     `json:"-"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
