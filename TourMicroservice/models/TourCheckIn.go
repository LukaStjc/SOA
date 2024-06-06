package models

import (
	"time"

	"gorm.io/gorm"
)

type TourCheckIn struct {
	gorm.Model
	ID           uint      `gorm:"primaryKey"`
	KeyPointID   uint      `json:"keyPointID" gorm:"type:uint;foreignKey:KeyPointID"`
	UserID       uint      `json:"userID" gorm:"type:uint;foreignKey:UserID"`
	VisitingTime time.Time `json:"visitingTime"`
}
