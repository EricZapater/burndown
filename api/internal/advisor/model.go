package advisor

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type AnalysisJob struct {
	ID          uuid.UUID `json:"id"`
	Status      string    `json:"status"`
	ImagePath   string    `json:"image_path"`
	Description string    `json:"description"`
	RequestMeta json.RawMessage `json:"request_meta"`
	Response json.RawMessage `json:"response"`
	ErrorMessage          string          `json:"error_message"`
	RequestedAt  time.Time       `json:"requested_at"`
    ResponsedAt  *time.Time      `json:"responsed_at"`
}