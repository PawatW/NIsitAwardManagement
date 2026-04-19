package models

import "github.com/google/uuid"

type Faculty struct {
	ID   uuid.UUID `db:"id" gorm:"primaryKey;type:uuid"`
	Name string    `db:"name"`
	CampusID    uuid.UUID    `db:"campus_id" gorm:"type:uuid;not null"`
	Campus      Campus       `gorm:"foreignKey:CampusID"`
	Departments []Department `gorm:"foreignKey:FacultyID"`
}