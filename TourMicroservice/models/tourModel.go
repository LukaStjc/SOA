package models

import (
	"gorm.io/gorm"
)

type TourType int

const (
	Easy TourType = iota + 1
	Moderate
	Hard
)

type Tour struct {
	gorm.Model
	Name        string     `json:"name" gorm:"not null;type:string"`
	Description string     `json:"description" gorm:"not null;type:string"`
	Type        TourType   `json:"type"`
	Tags        string     `json:"tags" gorm:"type:string"`
	Price       float64    `json:"price"`
	UserID      uint       `json:"userId" gorm:"type:uint;foreignKey:UserID"`
	KeyPoints   []KeyPoint `json:"keyPoints" gorm:"foreignKey:TourID"`
}
