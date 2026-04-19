package models

import "time"

//academic term represents a semester/term for award applications
type AcademicTerm struct {
	ID           int       `gorm:"primaryKey;autoIncrement" json:"id"`
	AcademicYear int       `gorm:"not null" json:"academic_year"`                    // e.g., 2025
	Semester     string    `gorm:"type:varchar(20);not null" json:"semester"`        // "first" or "second"
	StartDate    time.Time `gorm:"type:date;not null" json:"start_date"`
	EndDate      time.Time `gorm:"type:date;not null" json:"end_date"`
	IsOpen       bool      `gorm:"default:false;not null" json:"is_open"`            // Can students apply?
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	
	//requests that belong to this term
	Requests []Request `gorm:"foreignKey:AcademicTermID" json:"requests,omitempty"`
}

func (AcademicTerm) TableName() string {
	return "academic_terms"
}
