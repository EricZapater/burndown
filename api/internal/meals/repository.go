package meals

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Repository interface {
	Create(ctx context.Context, meal *Meal) error
	FindAll(ctx context.Context) ([]Meal, error)
	FindAllByUserId(ctx context.Context, userId string) ([]Meal, error)
	FindAllByUserIdAndDate(ctx context.Context, userId string, date string) ([]Meal, error)	
	FindById(ctx context.Context, id string) (*Meal, error)
	Update(ctx context.Context, meal *Meal) error
	Delete(ctx context.Context, id string) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, meal *Meal) error {
	sql := `INSERT INTO meals(id, user_id, timestamp, name, description, image_path, calories, protein_g, carbs_g, fat_g, confidence_level)
			VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err := r.db.ExecContext(ctx, sql, meal.ID, meal.UserID, meal.Timestamp, meal.Name, meal.Description, meal.ImagePath, meal.Calories, meal.Protein, meal.Carbs, meal.Fat, meal.Confidence)
	log.Println(err)
	return err
}

func (r *repository) FindAll(ctx context.Context) ([]Meal, error) {
	sql := `SELECT id, user_id, timestamp, name, description, image_path, calories, protein_g, carbs_g, fat_g, confidence_level
			FROM meals`
	rows, err := r.db.QueryContext(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var meals []Meal
	for rows.Next() {
		var meal Meal
		err := rows.Scan(&meal.ID, &meal.UserID, &meal.Timestamp, &meal.Name, &meal.Description, &meal.ImagePath, &meal.Calories, &meal.Protein, &meal.Carbs, &meal.Fat, &meal.Confidence)
		if err != nil {
			return nil, err
		}
		meals = append(meals, meal)
	}
	return meals, nil
}

func (r *repository) FindAllByUserId(ctx context.Context, userId string) ([]Meal, error) {
	sql := `SELECT id, user_id, timestamp, name, description, image_path, calories, protein_g, carbs_g, fat_g, confidence_level
			FROM meals
			WHERE user_id = $1`
	rows, err := r.db.QueryContext(ctx, sql, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var meals []Meal
	for rows.Next() {
		var meal Meal
		err := rows.Scan(&meal.ID, &meal.UserID, &meal.Timestamp, &meal.Name, &meal.Description, &meal.ImagePath, &meal.Calories, &meal.Protein, &meal.Carbs, &meal.Fat, &meal.Confidence)
		if err != nil {
			return nil, err
		}
		meals = append(meals, meal)
	}
	return meals, nil
}

func (r *repository) FindAllByUserIdAndDate(ctx context.Context, userId string, date string) ([]Meal, error) {
	layout := "2006-01-02"
    startDate, err := time.Parse(layout, date)
    if err != nil {
        return nil, fmt.Errorf("invalid date format (expected YYYY-MM-DD): %w", err)
    }
    
    endDate := startDate.Add(24 * time.Hour)
	sql := `SELECT id, user_id, timestamp, name, description, image_path, calories, protein_g, carbs_g, fat_g, confidence_level
			FROM meals
			WHERE user_id = $1 AND timestamp >= $2 AND timestamp < $3
			ORDER BY timestamp DESC`
	rows, err := r.db.QueryContext(ctx, sql, userId, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var meals []Meal
	for rows.Next() {
		var meal Meal
		err := rows.Scan(&meal.ID, &meal.UserID, &meal.Timestamp, &meal.Name, &meal.Description, &meal.ImagePath, &meal.Calories, &meal.Protein, &meal.Carbs, &meal.Fat, &meal.Confidence)
		if err != nil {
			return nil, err
		}
		meals = append(meals, meal)
	}
	return meals, nil
}

func (r *repository) FindById(ctx context.Context, id string) (*Meal, error) {
	sql := `SELECT id, user_id, timestamp, name, description, image_path, calories, protein_g, carbs_g, fat_g, confidence_level
			FROM meals
			WHERE id = $1`
	row := r.db.QueryRowContext(ctx, sql, id)
	var meal Meal
	err := row.Scan(&meal.ID, &meal.UserID, &meal.Timestamp, &meal.Name, &meal.Description, &meal.ImagePath, &meal.Calories, &meal.Protein, &meal.Carbs, &meal.Fat, &meal.Confidence)
	if err != nil {
		return nil, err
	}
	return &meal, nil
}

func (r *repository) Update(ctx context.Context, meal *Meal) error {
	sql := `UPDATE meals
			SET user_id = $2, timestamp = $3, name = $4, description = $5, image_path = $6, calories = $7, protein_g = $8, carbs_g = $9, fat_g = $10, confidence_level = $11
			WHERE id = $1`
	_, err := r.db.ExecContext(ctx, sql, meal.ID, meal.UserID, meal.Timestamp, meal.Name, meal.Description, meal.ImagePath, meal.Calories, meal.Protein, meal.Carbs, meal.Fat, meal.Confidence)
	return err
}

func (r *repository) Delete(ctx context.Context, id string) error {
	sql := `DELETE FROM meals
			WHERE id = $1`
	_, err := r.db.ExecContext(ctx, sql, id)
	return err
}
