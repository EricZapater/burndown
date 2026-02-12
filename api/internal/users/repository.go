package users

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, user *User) error
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

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, user *User) error {
	sql := `INSERT INTO users (id, name, email, password, active, is_admin) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, sql, user.ID, user.Name, user.Email, user.Password, user.Active, user.IsAdmin)
	return err
}

func (r *repository) GetAll(ctx context.Context) ([]User, error) {
	sql := `SELECT id, name, email, active, is_admin FROM users ORDER BY name ASC`
	rows, err := r.db.QueryContext(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Active, &user.IsAdmin); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *repository) ChangePassword(ctx context.Context, id uuid.UUID, password string) error {
	sql := `UPDATE users SET password = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, sql, password, id)
	return err
}
func (r *repository)InactiveUser(ctx context.Context, id uuid.UUID) error {
	sql := `UPDATE users SET active = false WHERE id = $1`
	_, err := r.db.ExecContext(ctx, sql, id)
	return err
}

func (r *repository) FindByEmail(ctx context.Context, email string) (*User, error) {
	sql := `SELECT id, name, email, password, active, is_admin FROM users WHERE email = $1`
	row := r.db.QueryRowContext(ctx, sql, email)
	user := &User{}
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Active, &user.IsAdmin)
	return user, err
}

func (r *repository) FindByID(ctx context.Context, id uuid.UUID) (*User, error) {
	sql := `SELECT id, name, email, password, active, is_admin FROM users WHERE id = $1`
	row := r.db.QueryRowContext(ctx, sql, id)
	user := &User{}
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Active, &user.IsAdmin)
	return user, err
}




func (r *repository) AddMacros(ctx context.Context, id uuid.UUID, macros Macros) error {
	sql := `INSERT INTO macros (id, user_id, timestamp, calories, protein_g, carbs_g, fat_g, active) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.ExecContext(ctx, sql, macros.ID, id, macros.Timestamp, macros.Calories, macros.Protein, macros.Carbs, macros.Fat, macros.Active)
	return err
}

func (r *repository) GetMacros(ctx context.Context, id uuid.UUID) ([]Macros, error) {
	sql := `SELECT id, user_id, timestamp, calories, protein_g, carbs_g, fat_g, active FROM macros WHERE user_id = $1 AND active = true`
	rows, err := r.db.QueryContext(ctx, sql, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var macros []Macros
	for rows.Next() {
		macro := Macros{}
		err := rows.Scan(&macro.ID, &macro.UserID, &macro.Timestamp, &macro.Calories, &macro.Protein, &macro.Carbs, &macro.Fat, &macro.Active)
		if err != nil {
			return nil, err
		}
		macros = append(macros, macro)
	}
	return macros, nil
}

func (r *repository) GetHistoricalMacros(ctx context.Context, id uuid.UUID) ([]Macros, error) {
	sql := `SELECT id, user_id, timestamp, calories, protein_g, carbs_g, fat_g, active FROM macros WHERE user_id = $1 ORDER BY timestamp DESC`
	rows, err := r.db.QueryContext(ctx, sql, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var macros []Macros
	for rows.Next() {
		macro := Macros{}
		err := rows.Scan(&macro.ID, &macro.UserID, &macro.Timestamp, &macro.Calories, &macro.Protein, &macro.Carbs, &macro.Fat, &macro.Active)
		if err != nil {
			return nil, err
		}
		macros = append(macros, macro)
	}
	return macros, nil
}

func (r *repository) DeleteMacros(ctx context.Context, id uuid.UUID) error {
	sql := `DELETE FROM macros WHERE id = $1`
	_, err := r.db.ExecContext(ctx, sql, id)
	return err
}

func (r *repository) UpdateMacros(ctx context.Context, id uuid.UUID, macros Macros) error {
	sql := `UPDATE macros SET calories = $1, protein_g = $2, carbs_g = $3, fat_g = $4 WHERE id = $5`
	_, err := r.db.ExecContext(ctx, sql, macros.Calories, macros.Protein, macros.Carbs, macros.Fat, id)
	return err
}

func (r *repository) InactiveMacros(ctx context.Context, id uuid.UUID) error {
	sql := `UPDATE macros SET active = false WHERE id = $1`
	_, err := r.db.ExecContext(ctx, sql, id)
	return err
}