package models

import "time"

type User struct {
	ID        uint      `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Password  string    `json:"-" db:"password"`
	Role      string    `json:"role" db:"role"`
	CreatedAt time.Time `json:"-" db:"created_at"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}
