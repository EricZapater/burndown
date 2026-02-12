package users

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Active   bool   `json:"active"`
	IsAdmin bool `json:"is_admin"`
	Macros   []Macros `json:"macros"`
}

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Macros struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Timestamp time.Time `json:"timestamp"`
	Calories  int       `json:"calories"`
	Protein   float64   `json:"protein_g"`
	Carbs     float64   `json:"carbs_g"`
	Fat       float64   `json:"fat_g"`
	Active    bool      `json:"active"`
}