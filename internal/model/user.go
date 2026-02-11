package model

import "time"

type User struct {
	ID           int64     `json:"id" db:"id"`
	Name         string    `json:"name" db:"name" validate:"required,min=3,max=50"`
	Email        string    `json:"email" db:"email" validate:"required,email"`
	PasswordHash string    `json:"password_hash" db:"password_hash" validate:"required,min=6"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}
