package requests

// StudentRegisterRequest is for student self-registration via Google OAuth
// Other roles (ADMIN, DEAN, etc.) are created by Admin via Admin endpoints
// After registration, student must login again to get access token
// Campus, Faculty and Department are required for approval workflow
type StudentRegisterRequest struct {
	RegisterToken string `json:"register_token" validate:"required"`
	FirstName     string `json:"first_name" validate:"required"`
	LastName      string `json:"last_name" validate:"required"`
	PhoneNumber   string `json:"phone_number" validate:"required"`
	NisitID       string `json:"nisit_id" validate:"required"`
	CampusID      string `json:"campus_id" validate:"required,uuid4"`
	FacultyID     string `json:"faculty_id" validate:"required,uuid4"`
	DepartmentID  string `json:"department_id" validate:"required,uuid4"`
}
