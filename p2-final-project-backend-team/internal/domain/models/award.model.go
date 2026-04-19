package models

// This model represents award categories.
import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

// formField represents a single field in the dynamic form
type FormField struct {
	ID         string           `json:"id"`
	Label      string           `json:"label"`
	Type       string           `json:"type"`
	Required   bool             `json:"required"`
	Options    []string         `json:"options,omitempty"`
	Validation *FieldValidation `json:"validation,omitempty"`
}

type FieldValidation struct {
	MinLength *int    `json:"min_length,omitempty"`
	MaxLength *int    `json:"max_length,omitempty"`
	Min       *int    `json:"min,omitempty"`
	Max       *int    `json:"max,omitempty"`
	Pattern   *string `json:"pattern,omitempty"`
}

// FormStructure represents the dynamic form structure as a slice of FormField
type FormStructure []FormField

// Scan implements the sql.Scanner interface to scan a value from the database into the FormStructure
func (fs *FormStructure) Scan(value interface{}) error {

	if value == nil {
		*fs = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan FormStructure: type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, fs)
}

// value implements the driver.Valuer interface to convert the FormStructure to a value that can be stored in the database
func (fs FormStructure) Value() (driver.Value, error) {
	if fs == nil {
		return nil, nil
	}
	return json.Marshal(fs)
}

type AwardCategory struct {
	ID            int            `gorm:"primaryKey;autoIncrement" json:"id"`
	Name          string         `gorm:"type:varchar(255);not null" json:"name"`
	Description   string         `gorm:"type:text" json:"description"`
	FormStructure FormStructure  `gorm:"type:jsonb" json:"form_structure"`
	IsActive      bool           `gorm:"default:true" json:"is_active"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// what does this do? -> This function overrides the default table name for the AwardCategory model in the database.
// When  save or read AwardCategory, use the table named award_categories
func (AwardCategory) TableName() string {
	return "award_categories"
}
