package models

import (
	"time"
)

// user represents an application user

// type User struct {
// 	gorm.Model
// 	Username string `gorm:"not null" json:"username" binding:"required"`
// 	Email    string `gorm:"unique;not null" json:"email" binding:"required,email"`
// 	Password string `gorm:"not null" binding:"required"`
// 	Books    []Book `gorm:"foreignKey:OwnerID;references:ID"`
// }

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"not null" json:"username" binding:"required"`
	Email     string    `gorm:"unique;not null" json:"email" binding:"required,email"`
	Password  string    `gorm:"not null" binding:"required"`
	Books     []Book    `gorm:"foreignKey:OwnerID;references:ID;constraint:OnDelete:CASCADE" json:"books"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
