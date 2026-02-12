package users

import (
	"context"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Create(ctx context.Context, request *CreateUserRequest) error
	GetAll(ctx context.Context) ([]User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)
	ChangePassword(ctx context.Context, id uuid.UUID, password string) error
	InactiveUser(ctx context.Context, id uuid.UUID) error
	AddMacros(ctx context.Context, id uuid.UUID, macros Macros) error
	GetMacros(ctx context.Context, id uuid.UUID) ([]Macros, error)
	GetHistoricalMacros(ctx context.Context, id uuid.UUID) ([]Macros, error)
	DeleteMacros(ctx context.Context, id uuid.UUID) error
	UpdateMacros(ctx context.Context, id uuid.UUID, macros Macros) error
	InactiveMacros(ctx context.Context, id uuid.UUID) error
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{repository: repository}
}

func (s *service) Create(ctx context.Context, request *CreateUserRequest) error {
	password, err := s.encryptPassword(request.Password)
	if err != nil {
		return err
	}
	user := &User{
		ID:       uuid.New(),
		Name:     request.Name,
		Email:    request.Email,
		Password: password,
		Active:   true,
	}
	return s.repository.Create(ctx, user)
}

func (s *service) GetAll(ctx context.Context) ([]User, error) {
	return s.repository.GetAll(ctx)
}

func (s *service) FindByEmail(ctx context.Context, email string) (*User, error) {
	return s.repository.FindByEmail(ctx, email)
}

func (s *service) FindByID(ctx context.Context, id uuid.UUID) (*User, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *service) ChangePassword(ctx context.Context, id uuid.UUID, password string) error {
	encryptedPassword, err := s.encryptPassword(password)
	if err != nil {
		return err
	}
	return s.repository.ChangePassword(ctx, id, encryptedPassword)
}

func (s *service) InactiveUser(ctx context.Context, id uuid.UUID) error {
	return s.repository.InactiveUser(ctx, id)
}

func (s *service) AddMacros(ctx context.Context, id uuid.UUID, macros Macros) error {
	existentMacros,err := s.repository.GetMacros(ctx, id)
	if err != nil {
		return err
	}
	for _, macro := range existentMacros {
		if macro.Active {
			s.repository.InactiveMacros(ctx, macro.ID)
		}
	}
	macros.ID = uuid.New()
	macros.Active = true
	macros.Timestamp = time.Now()
	return s.repository.AddMacros(ctx, id, macros)
}

func (s *service) GetMacros(ctx context.Context, id uuid.UUID) ([]Macros, error) {
	return s.repository.GetMacros(ctx, id)
}

func (s *service) GetHistoricalMacros(ctx context.Context, id uuid.UUID) ([]Macros, error) {
	return s.repository.GetHistoricalMacros(ctx, id)
}

func (s *service) DeleteMacros(ctx context.Context, id uuid.UUID) error {
	return s.repository.DeleteMacros(ctx, id)
}

func (s *service) UpdateMacros(ctx context.Context, id uuid.UUID, macros Macros) error {
	return s.repository.UpdateMacros(ctx, id, macros)
}

func (s *service) InactiveMacros(ctx context.Context, id uuid.UUID) error {
	return s.repository.InactiveMacros(ctx, id)
}

func(s *service)encryptPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }	
	return string(hashedPassword), nil
}