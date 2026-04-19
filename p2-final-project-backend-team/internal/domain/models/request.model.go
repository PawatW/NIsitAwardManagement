package models

import (
    "time"
    "database/sql/driver"
    "encoding/json"
    "errors"
    "github.com/google/uuid"
)

// AdditionalData stores dynamic form data as JSON
type AdditionalData map[string]interface{}

// Scan implements sql.Scanner for reading from database
func (ad *AdditionalData) Scan(value interface{}) error {
    if value == nil {
        *ad = AdditionalData{}
        return nil
    }
    
    bytes, ok := value.([]byte)
    if !ok {
        return errors.New("failed to scan AdditionalData")
    }
    
    return json.Unmarshal(bytes, ad)
}

// Value implements driver.Valuer for writing to database
func (ad AdditionalData) Value() (driver.Value, error) {
    if len(ad) == 0 {
        return nil, nil
    }
    return json.Marshal(ad)
}

type Request struct {
    ID             uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
    StudentID      uuid.UUID `gorm:"type:uuid;not null" json:"student_id"`
    AwardID        int       `gorm:"not null" json:"award_id"`
    AcademicTermID int       `gorm:"not null" json:"academic_term_id"`

    // Dynamic form data based on award category's form_structure
    AdditionalData AdditionalData `gorm:"type:jsonb" json:"additional_data"`

    // Snapshot data for PDF generation
    SnapshotStudyYear   int      `gorm:"not null" json:"snapshot_study_year"`
    SnapshotGPA         float64  `gorm:"type:decimal(3,2);not null" json:"snapshot_gpa"`
    SnapshotAdvisor     string   `gorm:"type:varchar(255);not null" json:"snapshot_advisor"`
    SnapshotDateOfBirth *string  `gorm:"type:date" json:"snapshot_date_of_birth,omitempty"`
    SnapshotPhone       *string  `gorm:"type:varchar(50)" json:"snapshot_phone,omitempty"`
    SnapshotAddress     *string  `gorm:"type:text" json:"snapshot_address,omitempty"`

    CurrentStatus string    `gorm:"type:varchar(50);not null" json:"current_status"`
    CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`

    // Foreign key relations
    Student       User                   `gorm:"foreignKey:StudentID;references:UserID" json:"student,omitempty"`
    Award         AwardCategory          `gorm:"foreignKey:AwardID;references:ID" json:"award,omitempty"`
    AcademicTerm  AcademicTerm           `gorm:"foreignKey:AcademicTermID;references:ID" json:"academic_term,omitempty"`
    
    // One-to-many relations
    StatusHistory   []RequestStatusHistory `gorm:"foreignKey:RequestID" json:"status_history,omitempty"`
    Documents       []RequestDocument      `gorm:"foreignKey:RequestID" json:"documents,omitempty"`
    CategoryChanges []AwardCategoryChange  `gorm:"foreignKey:RequestID" json:"category_changes,omitempty"`
}

func (Request) TableName() string {
    return "requests"
}