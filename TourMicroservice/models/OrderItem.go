package models

import (
	"gorm.io/gorm"
)

type OrderItem struct {
	gorm.Model
	TourId         uint    `json:"tourId" gorm:"type:uint;foreignKey:TourID"`
	TourName       string  `json:"tourName" gorm:"type:string"`
	TourPrice      float64 `json:"tourPrice"`
	NumberOfPeople int     `json:"numberOfPeople"`
}
