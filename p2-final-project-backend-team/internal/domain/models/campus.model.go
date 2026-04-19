package models

import (
	"github.com/google/uuid"
)


type Campus struct {
	ID uuid.UUID `db:"id" gorm:"primaryKey;type:uuid"`
	Name string `db:"name"`
	Faculties []Faculty `gorm:"foreignKey:CampusID"`
}