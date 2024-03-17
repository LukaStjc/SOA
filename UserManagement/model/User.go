package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRole predstavlja ulogu korisnika.
type UserRole int

// Definisanje konstanti za uloge korisnika.
const (
	Guide UserRole = iota
	Tourist
	Administrator
)

// Funkcija String vraća string reprezentaciju uloge korisnika.
func (role UserRole) String() string {
	names := [...]string{"Guide", "Tourist", "Administrator"}
	if role < Guide || role > Administrator {
		return "Unknown"
	}
	return names[role]
}

// Struktura User predstavlja korisnika.
type User struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username" gorm:"not null;type:string"`
	Password string    `json:"password" gorm:"not null;type:string"`
	Email    string    `json:"email" gorm:"not null;type:string"`
	Role     UserRole  `json:"role"`
}

func (user *User) BeforeCreate(scope *gorm.DB) error {
	user.ID = uuid.New()
	return nil
}
