package models

import (
	"time"
	"github.com/google/uuid"
)

type RequestStatusHistory struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	RequestID uuid.UUID `gorm:"type:uuid;not null" json:"request_id"`
	Status    string    `gorm:"type:varchar(50);not null" json:"status"`
	UpdatedBy uuid.UUID `gorm:"type:uuid;not null" json:"updated_by"`
	Remark    string    `gorm:"type:text" json:"remark"`
	UpdatedAt time.Time `gorm:"autoCreateTime" json:"updated_at"`
	
	//references:UserID
	Request Request `gorm:"foreignKey:RequestID;references:ID" json:"request,omitempty"`
	User    User    `gorm:"foreignKey:UpdatedBy;references:UserID" json:"updated_by_user,omitempty"`
}

func (RequestStatusHistory) TableName() string {
	return "request_status_history"
}
