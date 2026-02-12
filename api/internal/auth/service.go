package auth

import (
	"context"
	"errors"

	jwt "github.com/appleboy/gin-jwt/v2"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Login(ctx context.Context, request *LoginRequest) (*LoginResponse, error)
}

type service struct {
	userProvider UserProvider
	jwtMiddleware *jwt.GinJWTMiddleware
}

func NewService(userProvider UserProvider, jwtMiddleware *jwt.GinJWTMiddleware) Service {
	return &service{userProvider: userProvider, jwtMiddleware: jwtMiddleware}
}

func (s *service) Login(ctx context.Context, request *LoginRequest) (*LoginResponse, error) {
	user, err := s.userProvider.FindByEmail(ctx, request.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		return nil, errors.New("invalid password")
	}
	token, expire, err := s.jwtMiddleware.TokenGenerator(user.ID.String())
	if err != nil {
		return nil, err
	}
	return &LoginResponse{
		Token: token,
		User:  *user,
		Expire: expire,
	}, nil
}
