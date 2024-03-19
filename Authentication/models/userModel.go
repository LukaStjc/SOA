package models

type User struct {
	ID       uint   `gorm:"primaryKey"` // Explicitly declare the ID field as the primary key
	Username string `gorm:"unique"`
	Password string
	Role     uint8
}
