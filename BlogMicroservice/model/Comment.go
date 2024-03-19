package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Comment struct {
	ID             uuid.UUID `json:"id"`
	UserID         uuid.UUID `json:"user id" gorm:"type:uuid;foreignKey:UserID"`
	BlogID         uuid.UUID `json:"blog id" gorm:"type:uuid;foreignKey:BlogID"`
	PublishDate    time.Time `json:"publish date" gorm:"not null;type:timestamp"`
	Text           string    `json:"text" gorm:"not null;type:text"`
	LastChangeDate time.Time `json:"last change date" gorm:"not null;type:timestamp"`
}

func (comment *Comment) BeforeCreate(scope *gorm.DB) error {
	comment.ID = uuid.New()
	return nil
}
