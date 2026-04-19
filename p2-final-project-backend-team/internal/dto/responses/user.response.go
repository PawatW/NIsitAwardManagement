package responses

import (
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/domain/models"
	"github.com/google/uuid"
)

// UserResponse - Basic user information for responses
type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	NisitID   string    `json:"nisit_id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`

	// Optional: Include related data if needed
	CampusID  *uuid.UUID `json:"campus_id,omitempty"`
	FacultyID *uuid.UUID `json:"faculty_id,omitempty"`
	MajorID   *uuid.UUID `json:"major_id,omitempty"`
}

// UserBasicResponse - Minimal user info (for nested responses)
type UserBasicResponse struct {
	ID             uuid.UUID `json:"id"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	NisitID        string    `json:"nisit_id"`
	PhoneNumber    *string   `json:"phone_number,omitempty"`
	CampusName     string    `json:"campus_name,omitempty"`
	FacultyName    string    `json:"faculty_name,omitempty"`
	DepartmentName string    `json:"department_name,omitempty"`
}

// OrganizationResponse - Organization info for nested responses
type OrganizationResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// FullUserResponse - Complete user info with all related data
type FullUserResponse struct {
	ID           uuid.UUID             `json:"id"`
	FirstName    string                `json:"first_name"`
	LastName     string                `json:"last_name"`
	Email        string                `json:"email"`
	PhoneNumber  *string               `json:"phone_number,omitempty"`
	ProfileURL   *string               `json:"profile_url,omitempty"`
	Role         string                `json:"role"`
	NisitID      *string               `json:"nisit_id,omitempty"`
	AuthProvider string                `json:"auth_provider"`
	IsActive     bool                  `json:"is_active"`
	Campus       *OrganizationResponse `json:"campus,omitempty"`
	Faculty      *OrganizationResponse `json:"faculty,omitempty"`
	Department   *OrganizationResponse `json:"department,omitempty"`
}

// ToFullUserResponse converts User model to FullUserResponse
func ToFullUserResponse(u *models.User) *FullUserResponse {
	if u == nil {
		return nil
	}

	resp := &FullUserResponse{
		ID:           u.UserID,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		Email:        u.Email,
		PhoneNumber:  u.PhoneNumber,
		ProfileURL:   u.ProfileURL,
		Role:         string(u.Role),
		NisitID:      u.NisitID,
		AuthProvider: u.AuthProvider,
		IsActive:     u.IsActive,
	}

	if u.Campus != nil {
		resp.Campus = &OrganizationResponse{
			ID:   u.Campus.ID,
			Name: u.Campus.Name,
		}
	}

	if u.Faculty != nil {
		resp.Faculty = &OrganizationResponse{
			ID:   u.Faculty.ID,
			Name: u.Faculty.Name,
		}
	}

	if u.Department != nil {
		resp.Department = &OrganizationResponse{
			ID:   u.Department.ID,
			Name: u.Department.Name,
		}
	}

	return resp
}
