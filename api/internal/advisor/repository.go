package advisor

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type Repository interface {
	CreateAnalysisJob(ctx context.Context, job *AnalysisJob) error
	GetAnalysisJob(ctx context.Context, id uuid.UUID) (*AnalysisJob, error)
	UpdateAnalysisJob(ctx context.Context, job *AnalysisJob) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateAnalysisJob(ctx context.Context, job *AnalysisJob) error {
	sql := `INSERT INTO analysis_jobs (id, status, image_path, description, request_meta, response, error_message, requested_at, responsed_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.db.ExecContext(ctx, sql, job.ID, job.Status, job.ImagePath, job.Description, job.RequestMeta, job.Response, job.ErrorMessage, job.RequestedAt, job.ResponsedAt)
	return err
}

func (r *repository) GetAnalysisJob(ctx context.Context, id uuid.UUID) (*AnalysisJob, error) {
	sql := `SELECT id, status, image_path, description, request_meta, response, error_message, requested_at, responsed_at
	FROM analysis_jobs
	WHERE id = $1`
	row := r.db.QueryRowContext(ctx, sql, id)
	var job AnalysisJob
	err := row.Scan(&job.ID, &job.Status, &job.ImagePath, &job.Description, &job.RequestMeta, &job.Response, &job.ErrorMessage, &job.RequestedAt, &job.ResponsedAt)
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func (r *repository) UpdateAnalysisJob(ctx context.Context, job *AnalysisJob) error {
	sql := `UPDATE analysis_jobs
	SET status = $2, image_path = $3, description = $4, request_meta = $5, response = $6, error_message = $7, requested_at = $8, responsed_at = $9
	WHERE id = $1`
	_, err := r.db.ExecContext(ctx, sql, job.ID, job.Status, job.ImagePath, job.Description, job.RequestMeta, job.Response, job.ErrorMessage, job.RequestedAt, job.ResponsedAt)
	return err
}