package responses

import (
	"time"
	"github.com/google/uuid"
)

type AwardCategoryChangeResponse struct {
	ID            int       `json:"id"`
	RequestID     uuid.UUID `json:"request_id"`
	OldCategoryID int       `json:"old_category_id"`
	NewCategoryID int       `json:"new_category_id"`
	ChangedAt     time.Time `json:"changed_at"`
	ChangedBy     uuid.UUID `json:"changed_by"`
	Remark        string    `json:"remark,omitempty"`
	
	//optional  Include category names
	OldCategory   *AwardCategoryResponse `json:"old_category,omitempty"`
	NewCategory   *AwardCategoryResponse `json:"new_category,omitempty"`
	
	//wh chnahed
	ChangedByUser *UserBasicResponse `json:"changed_by_user,omitempty"`
}
