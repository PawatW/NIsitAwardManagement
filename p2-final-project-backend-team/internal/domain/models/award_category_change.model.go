package models

import (
	"time"
	"github.com/google/uuid"
)

type AwardCategoryChange struct {
	ID            int       `gorm:"primaryKey;autoIncrement" json:"id"`
	RequestID     uuid.UUID `gorm:"type:uuid;not null" json:"request_id"`
	OldCategoryID int       `gorm:"not null" json:"old_category_id"`
	NewCategoryID int       `gorm:"not null" json:"new_category_id"`
	ChangedBy     uuid.UUID `gorm:"type:uuid;not null" json:"changed_by"`
	ChangedAt     time.Time `gorm:"autoCreateTime" json:"changed_at"`
	
	//refer:UserID
	Request     Request       `gorm:"foreignKey:RequestID;references:ID" json:"request,omitempty"`
	OldCategory AwardCategory `gorm:"foreignKey:OldCategoryID;references:ID" json:"old_category,omitempty"`
	NewCategory AwardCategory `gorm:"foreignKey:NewCategoryID;references:ID" json:"new_category,omitempty"`
	User        User          `gorm:"foreignKey:ChangedBy;references:UserID" json:"changed_by_user,omitempty"`
}

func (AwardCategoryChange) TableName() string {
	return "award_category_changes"
}
