package auth

import (
	"api/internal/users"
	"context"
)

type UserProvider interface {
	FindByEmail(ctx context.Context, email string) (*users.User, error)
}