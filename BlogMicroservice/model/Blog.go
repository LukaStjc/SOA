package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BlogStatus int

const (
	Draft     BlogStatus = iota // 0
	Published                   // 1
	Closed                      // 2
)

type Blog struct {
	ID          uuid.UUID  `json:"id"`
	UserID      int        `json:"user id" gorm:"type:int;foreignKey:UserID"`
	Title       string     `json:"title" gorm:"not null;type:string"`
	Description string     `json:"description" gorm:"not null;type:text"`
	PublishDate time.Time  `json:"publish date" gorm:"not null;type:timestamp"`
	Status      BlogStatus `json:"status gorm:type:int"`
}

func (blog *Blog) BeforeCreate(scope *gorm.DB) error {
	blog.ID = uuid.New()
	blog.PublishDate = time.Now()
	return nil
}

func (status BlogStatus) StatusToString() string {
	switch status {
	case 0:
		return "Draft"
	case 1:
		return "Published"
	case 2:
		return "Closed"
	default:
		return "Undefined"
	}

}
