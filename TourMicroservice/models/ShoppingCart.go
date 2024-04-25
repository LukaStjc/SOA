package models

import (
	"gorm.io/gorm"
)

type ShoppingCart struct {
	gorm.Model
	OrderItems []uint  `json:"orderItems" gorm:"type:uint;foreignKey:TourID"`
	UserID     uint    `json:"userId" gorm:"type:uint;foreignKey:UserID"`
	Price      float64 `json:"price"`
}
