package models

import (
	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/domain/enums"
	"github.com/google/uuid"
)

type User struct {
	UserID       uuid.UUID      `db:"user_id" gorm:"primaryKey;type:uuid"`
	FirstName    string         `db:"first_name"`
	LastName     string         `db:"last_name"`
	Email        string         `db:"email" gorm:"unique"`
	PhoneNumber  *string        `db:"phone_number" gorm:"type:varchar(20)"`
	ProfileURL   *string        `db:"profile_url"`
	AuthProvider string         `db:"auth_provider" gorm:"default:'google'"`
	Role         enums.UserRole `db:"role" gorm:"not null"`
	NisitID      *string        `db:"nisit_id" gorm:"type:varchar(10);unique"`
	IsActive     bool           `db:"is_active" gorm:"not null;default:true" json:"is_active"`

	// Organizational references
	CampusID     *uuid.UUID `db:"campus_id" gorm:"type:uuid"`
	FacultyID    *uuid.UUID `db:"faculty_id" gorm:"type:uuid"`
	DepartmentID *uuid.UUID `db:"department_id" gorm:"type:uuid"`

	// Relations (for preloading)
	Campus     *Campus     `gorm:"foreignKey:CampusID"`
	Faculty    *Faculty    `gorm:"foreignKey:FacultyID"`
	Department *Department `gorm:"foreignKey:DepartmentID"`
}
