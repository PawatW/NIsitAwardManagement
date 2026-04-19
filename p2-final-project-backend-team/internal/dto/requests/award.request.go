package requests

import (
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/domain/models"
	"github.com/google/uuid"
)

type CreateAwardCategoryRequest struct {
	Name          string               `json:"name" validate:"required"`
	Description   string               `json:"description"`
	FormStructure models.FormStructure `json:"form_structure" validate:"required,min=1"` //
	IsActive      *bool                `json:"is_active"`
}

type UpdateAwardCategoryRequest struct {
	Name          *string               `json:"name" validate:"required"`
	Description   *string               `json:"description"` // Change to pointer
	FormStructure *models.FormStructure `json:"form_structure" validate:"omitempty,min=1"`
	IsActive      *bool                 `json:"is_active"`
}

// HonorRollQueryDTO - Query parameters for honor roll display
type HonorRollQueryDTO struct {
	AcademicYear int       `query:"academic_year"` // Filter by academic year (e.g., 2567)
	Semester     string    `query:"semester"`      // Filter by semester ("first" or "second")
	CampusID     uuid.UUID `query:"campus_id"`     // Filter by campus UUID
	AwardID      int       `query:"award_id"`      // Filter by specific award category
}
