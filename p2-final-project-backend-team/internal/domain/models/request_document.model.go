package models

import (
	"time"

	"github.com/google/uuid"
)

type RequestDocument struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	RequestID    uuid.UUID `gorm:"type:uuid;not null" json:"request_id"`
	Bucket       string    `gorm:"type:varchar(255);not null" json:"bucket"`
	ObjectKey    string    `gorm:"type:varchar(500);not null" json:"object_key"`
	OriginalName string    `gorm:"type:varchar(255);not null" json:"original_name"`
	ContentType  string    `gorm:"type:varchar(100);not null" json:"content_type"`
	Size         int64     `gorm:"not null" json:"size"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`

	//fk
	Request Request `gorm:"foreignKey:RequestID;references:ID" json:"request,omitempty"`
}

func (RequestDocument) TableName() string {
	return "request_documents"
}
