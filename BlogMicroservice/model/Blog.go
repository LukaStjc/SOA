package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BlogStatus int

const (
	Draft BlogStatus = iota
	Published
	Closed
)

type Blog struct {
	Id          uuid.UUID  `json:"id"`
	Title       string     `json:"title" gorm:"not null;type:string"`
	Description string     `json:"description" gorm:"not null;type:text"`
	PublishDate time.Time  `json:"publish date" gorm:"not null;type:timestamp"`
	Status      BlogStatus `json:"status gorm:type:int"`
}

func (blog *Blog) BeforeCreate(scope *gorm.DB) error {
	blog.Id = uuid.New()
	return nil
}
