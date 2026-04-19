package responses

import (
	"time"

	"github.com/google/uuid"
)

type RequestDocumentResponse struct {
	ID           uuid.UUID `json:"id"`
	RequestID    uuid.UUID `json:"request_id"`
	Bucket       string    `json:"bucket"`
	ObjectKey    string    `json:"object_key"`
	OriginalName string    `json:"original_name"`
	ContentType  string    `json:"content_type"`
	Size         int64     `json:"size"`
	FileURL      string    `json:"file_url"`
	UploadedAt   time.Time `json:"uploaded_at"`
}
