package models

import (
	"gorm.io/gorm"
)

type OrderItem struct {
	gorm.Model
	ID             uint    `gorm:"primaryKey"`
	TourID         uint    `json:"tourID" gorm:"type:uint;foreignKey:TourID"`
	TourRate       uint    `json:"tourRate"`
	CartID         uint    `json:"cartID" gorm:"not null;foreignKey:CartID"`
	TourName       string  `json:"tourName" gorm:"type:string"`
	TourPrice      float64 `json:"tourPrice"`
	NumberOfPeople int     `json:"numberOfPeople"`
}
