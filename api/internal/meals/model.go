package meals

import (
	"time"

	"github.com/google/uuid"
)

type Meal struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	Timestamp  time.Time `json:"timestamp"`
	Name       string  `json:"name"`
	Description string `json:"description"`
	ImagePath string `json:"image_path"`
	Calories   int     `json:"calories"`
	Protein    float64 `json:"protein_g"`
	Carbs      float64 `json:"carbs_g"`
	Fat        float64 `json:"fat_g"`
	Confidence float64  `json:"confidence_level"`
}

type CreateMealRequest struct {
	UserID     string `json:"user_id" binding:"required"`
	Name       string  `json:"name" binding:"required"`
	Description string `json:"description"`
	ImagePath string `json:"image_path" binding:"required"`
	Calories   int     `json:"calories"`
	Protein    float64 `json:"protein_g"`
	Carbs      float64 `json:"carbs_g"`
	Fat        float64 `json:"fat_g"`
	Confidence float64  `json:"confidence_level"`
}