package models

import (
	"time"
)

// book represents data about a book

// type Book struct {
// 	gorm.Model
// 	Title      string `gorm:"not null" json:"title" binding:"required"`
// 	Edition    *uint  `json:"edition,omitempty"`
// 	Author     string `gorm:"not null" json:"author" binding:"required"`
// 	PagesCount uint   `gorm:"not null" json:"pages_count" binding:"required"`
// 	PagesRead  uint   `gorm:"not null" json:"pages_read" binding:"required"`
// 	OwnerID    uint   `json:"owner_id"`
// }

type Book struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Title      string    `gorm:"not null" json:"title" binding:"required"`
	Edition    *uint     `json:"edition,omitempty"`
	Author     string    `gorm:"not null" json:"author" binding:"required"`
	PagesCount uint      `gorm:"not null" json:"pages_count" binding:"required"`
	PagesRead  uint      `gorm:"not null" json:"pages_read" binding:"required"`
	OwnerID    uint      `json:"owner_id"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
