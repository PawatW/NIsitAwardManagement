package responses

import (
	"time"

	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/domain/models"
	"github.com/google/uuid"
)

type AwardCategoryResponse struct {
	ID            int                  `json:"id"`
	Name          string               `json:"name"`
	Description   string               `json:"description"`
	FormStructure models.FormStructure `json:"form_structure"`
	IsActive      bool                 `json:"is_active"`
	CreatedAt     time.Time            `json:"created_at"`
	UpdatedAt     time.Time            `json:"updated_at"`
}

// WinnerResponse - Individual winner information for honor roll
type WinnerResponse struct {
	StudentID      uuid.UUID `json:"student_id"`
	StudentName    string    `json:"student_name"` // FirstName + LastName
	NisitID        string    `json:"nisit_id"`
	FacultyName    string    `json:"faculty_name"`
	DepartmentName string    `json:"department_name"`
	CampusName     string    `json:"campus_name"`
	GPA            float64   `json:"gpa"`
}

// AwardCategoryGroupResponse - Winners grouped by award category
type AwardCategoryGroupResponse struct {
	AwardID   int              `json:"award_id"`
	AwardName string           `json:"award_name"`
	Winners   []WinnerResponse `json:"winners"`
}

// AcademicYearGroupResponse - Categories grouped by academic year
type AcademicYearGroupResponse struct {
	AcademicYear int                          `json:"academic_year"`
	Semester     string                       `json:"semester"`
	Categories   []AwardCategoryGroupResponse `json:"categories"`
}

// HonorRollResponse - Main response for honor roll endpoint
type HonorRollResponse struct {
	Data         []AcademicYearGroupResponse `json:"data"`
	TotalWinners int                         `json:"total_winners"`
}
