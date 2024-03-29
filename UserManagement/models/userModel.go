package models

import (
	"gorm.io/gorm"
)

// UserRole predstavlja ulogu korisnika.
type UserRole int

// Definisanje konstanti za uloge korisnika.
const (
	Guide UserRole = iota + 1
	Tourist
	Administrator
)

// Funkcija String vraća string reprezentaciju uloge korisnika.
func (role UserRole) String() string {
	names := [...]string{"Guide", "Tourist", "Administrator"}
	if role < Guide || role > Administrator {
		return "Unknown"
	}
	return names[role-1]
}

// Struktura User predstavlja korisnika.
type User struct {
	gorm.Model
	Username string   `json:"username" gorm:"unique;not null;type:string"`
	Password string   `json:"password" gorm:"not null;type:string"`
	Email    string   `json:"email" gorm:"unique;not null;type:string"`
	Role     UserRole `json:"role"`
	Blocked  bool     `json:"blocked" gorm:"default:false"`
	Follows  []*User  `gorm:"many2many:user_follows;"`
}

// func (user *User) BeforeCreate(scope *gorm.DB) error {
// 	user.ID = uuid.New()
// 	return nil
// }
