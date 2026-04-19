package responses

import (
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/domain/models"
	"github.com/google/uuid"
	"time"
)

// RequestResponse - Full request details
type RequestResponse struct {
	ID             uuid.UUID             `json:"id"`
	StudentID      uuid.UUID             `json:"student_id"`
	AwardID        int                   `json:"award_id"`
	AcademicTermID int                   `json:"academic_term_id"`
	AdditionalData models.AdditionalData `json:"additional_data"`

	// Snapshot data
	SnapshotStudyYear   int     `json:"snapshot_study_year,omitempty"`
	SnapshotGPA         float64 `json:"snapshot_gpa,omitempty"`
	SnapshotAdvisor     string  `json:"snapshot_advisor,omitempty"`
	SnapshotDateOfBirth *string `json:"snapshot_date_of_birth,omitempty"`
	SnapshotPhone       *string `json:"snapshot_phone,omitempty"`
	SnapshotAddress     *string `json:"snapshot_address,omitempty"`

	CurrentStatus string    `json:"current_status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// Related data (populated when needed)
	Student       *UserBasicResponse             `json:"student,omitempty"`
	Award         *AwardCategoryResponse         `json:"award,omitempty"`
	AcademicTerm  *AcademicTermResponse          `json:"academic_term,omitempty"`
	StatusHistory []RequestStatusHistoryResponse `json:"status_history,omitempty"`
	Documents     []RequestDocumentResponse      `json:"documents,omitempty"`
}

// RequestListResponse - Simplified for listing (without relations)
type RequestListResponse struct {
	ID            uuid.UUID          `json:"id"`
	Student       *UserBasicResponse `json:"student,omitempty"`
	AwardName     string             `json:"award_name"`
	CurrentStatus string             `json:"current_status"`
	SnapshotGPA   float64            `json:"snapshot_gpa"`
	CreatedAt     time.Time          `json:"created_at"`
}

// AcademicTermResponse - Academic term info
type AcademicTermResponse struct {
	ID           int       `json:"id"`
	AcademicYear int       `json:"academic_year"`
	Semester     string    `json:"semester"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	IsOpen       bool      `json:"is_open"`
}

// PaginatedRequestResponse - For paginated lists
type PaginatedRequestResponse struct {
	Data       []*RequestListResponse `json:"data"`
	TotalItems int64                  `json:"total_items"`
	Page       int                    `json:"page"`
	Limit      int                    `json:"limit"`
	TotalPages int                    `json:"total_pages"`
}
