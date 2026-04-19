package models

import "github.com/google/uuid"

type Department struct {
	ID   uuid.UUID `db:"id" gorm:"primaryKey;type:uuid"`
	Name string    `db:"name"`
	FacultyID uuid.UUID `db:"faculty_id" gorm:"type:uuid;not null"`
	Faculty   Faculty   `gorm:"foreignKey:FacultyID"`
}