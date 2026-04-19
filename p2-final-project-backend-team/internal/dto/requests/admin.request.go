package requests

import "github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/domain/enums"

type AddUserRequest struct {
	Email string         `json:"email" validate:"required,email"`
	Role  enums.UserRole `json:"role" validate:"required,userrole"`
}

// CreateStaffRequest is used by admin to create staff/faculty users
type CreateStaffRequest struct {
	FirstName    string         `json:"first_name" validate:"required"`
	LastName     string         `json:"last_name" validate:"required"`
	Email        string         `json:"email" validate:"required,email"`
	PhoneNumber  string         `json:"phone_number" validate:"omitempty"`
	Role         enums.UserRole `json:"role" validate:"required,userrole"`
	CampusID     string         `json:"campus_id" validate:"omitempty,uuid4"`
	FacultyID    string         `json:"faculty_id" validate:"omitempty,uuid4"`
	DepartmentID string         `json:"department_id" validate:"omitempty,uuid4"`
}

type CreateCampusRequest struct {
	Name string `json:"name" validate:"required"`
}

type CreateFacultyRequest struct {
	Name     string `json:"name" validate:"required"`
	CampusID string `json:"campus_id" validate:"required,uuid4"`
}

type CreateDepartmentRequest struct {
	Name      string `json:"name" validate:"required"`
	FacultyID string `json:"faculty_id" validate:"required,uuid4"`
}

type UpdateUserStatusRequest struct {
	IsActive *bool `json:"is_active" validate:"required"`
}
