package requests

type CreateAcademicTermRequest struct {
	AcademicYear int    `json:"academic_year" validate:"required"`
	Semester     string `json:"semester" validate:"required,oneof=first second summer"`
	StartDate    string `json:"start_date" validate:"required,datetime=2006-01-02"`
	EndDate      string `json:"end_date" validate:"required,datetime=2006-01-02"`
	IsOpen       bool   `json:"is_open"`
}

type UpdateAcademicTermRequest struct {
	IsOpen bool `json:"is_open"`
}