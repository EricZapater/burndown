package meals

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, request *CreateMealRequest) error
	FindAll(ctx context.Context) ([]Meal, error)
	FindAllByUserId(ctx context.Context, userId string) ([]Meal, error)
	FindAllByUserIdAndDate(ctx context.Context, userId string, date string) ([]Meal, error)
	FindById(ctx context.Context, id string) (*Meal, error)
	Update(ctx context.Context, meal *Meal) error
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, request *CreateMealRequest) error {
	meal := &Meal{
		ID:          uuid.New(),
		UserID:      uuid.MustParse(request.UserID),
		Timestamp:   time.Now(),
		Name:        request.Name,
		Description: request.Description,
		ImagePath:   request.ImagePath,
		Calories:    request.Calories,
		Protein:     request.Protein,
		Carbs:       request.Carbs,
		Fat:         request.Fat,
		Confidence:  request.Confidence,
	}
	return s.repo.Create(ctx, meal)
}

func (s *service) FindAll(ctx context.Context) ([]Meal, error) {
	return s.repo.FindAll(ctx)
}

func (s *service) FindAllByUserId(ctx context.Context, userId string) ([]Meal, error) {
	return s.repo.FindAllByUserId(ctx, userId)
}

func (s *service) FindAllByUserIdAndDate(ctx context.Context, userId string, date string) ([]Meal, error) {
	return s.repo.FindAllByUserIdAndDate(ctx, userId, date)
}

func (s *service) FindById(ctx context.Context, id string) (*Meal, error) {
	return s.repo.FindById(ctx, id)
}

func (s *service) Update(ctx context.Context, meal *Meal) error {
	return s.repo.Update(ctx, meal)
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
