package models

import (
	"gorm.io/gorm"
)

type ShoppingCart struct {
	gorm.Model
	ID         uint         `gorm:"primaryKey"`
	OrderItems []*OrderItem `json:"orderItems" gorm:"foreignKey:CartID"`
	UserID     uint         `json:"userId" gorm:"type:uint;foreignKey:UserID"`
	Price      float64      `json:"price"`
}
