package auth

import (
	"api/internal/users"
	"time"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string    `json:"token"`
	Expire time.Time `json:"expire"`
	User  users.User `json:"user"`
}
