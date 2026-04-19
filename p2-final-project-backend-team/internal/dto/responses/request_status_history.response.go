package responses

import (
	"time"
	"github.com/google/uuid"
)

type RequestStatusHistoryResponse struct {
	ID        int       `json:"id"`
	RequestID uuid.UUID `json:"request_id"`
	Status    string    `json:"status"`
	ChangedAt time.Time `json:"changed_at"`
	ChangedBy string    `json:"changed_by"`
	Remark    string    `json:"remark,omitempty"` //omit empty -> if empty dont include
}