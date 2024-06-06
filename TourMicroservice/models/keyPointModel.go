package models

import (
	"gorm.io/gorm"
)

type KeyPoint struct {
	gorm.Model
	ID        uint    `gorm:"primaryKey"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	TourID    uint    `json:"tourId" gorm:"type:uint;foreignKey:TourID"`
}
