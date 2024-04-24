package models

import (
	"gorm.io/gorm"
)

// type KeyPointType int

// const (
// 	Beggining KeyPointType = iota + 1
// 	Middle
// 	Ending
// )

type KeyPoint struct {
	gorm.Model
	Name      string  `json:"name"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	TourID    int     `json:"tourId" gorm:"type:int;foreignKey:TourID"`
}
