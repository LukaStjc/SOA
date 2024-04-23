package models

import (
	"strings"

	"gorm.io/gorm"
)

// UserRole predstavlja ulogu korisnika.
type TourType int

// Definisanje konstanti za uloge korisnika.
const (
	Easy TourType = iota + 1
	Moderate
	Hard
)

// Struktura User predstavlja korisnika.
type Tour struct {
	gorm.Model
	Name        string   `json:"name" gorm:"not null;type:string"`
	Description string   `json:"description" gorm:"not null;type:string"`
	Type        TourType `json:"type"`
	Tags        string   `json:"tags" gorm:"type:string"`
}

func (t *Tour) AddTag(tag string) {
	if t.Tags != "" {
		t.Tags += ";" + tag
	} else {
		t.Tags = tag
	}
}

func (t *Tour) RemoveTag(tag string) {
	tags := strings.Split(t.Tags, ";")
	for i, existingTag := range tags {
		if existingTag == tag {
			tags = append(tags[:i], tags[i+1:]...)
			break
		}
	}
	t.Tags = strings.Join(tags, ";")
}

func (t *Tour) GetTags() []string {
	return strings.Split(t.Tags, ";")
}
