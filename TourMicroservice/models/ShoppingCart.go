package models

import (
	"gorm.io/gorm"
)

type ShoppingCart struct {
	gorm.Model
	OrderItems []*OrderItem `gorm:"many2many:orderItem_in_shoppingCart;"`
	UserID     uint         `json:"userId" gorm:"type:uint;foreignKey:UserID"`
	Price      float64      `json:"price"`
}
