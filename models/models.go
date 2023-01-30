package models

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	ID        uint
	Title     string
	Text      []byte
	Public    bool
	Slug      string
	CreatedAt time.Time
	UpdatedAt time.Time
	gorm.Model
}
type Page struct {
	ID        uint
	Title     string
	ShortName string
	Text      []byte
	CreatedAt time.Time
	UpdatedAt time.Time
	Sidebar   bool
	gorm.Model
}
