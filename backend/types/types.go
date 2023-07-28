package types

import (
	"errors"
	"regexp"
	"time"
)

func ValidateEmail(email string) bool {
	// Regular expression to match email format.
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

type Credentials struct {
	Email    string `form:"username" binding:"required,email"`
	Password string `form:"password" binding:"required"`
}

func (c *Credentials) ValidateCredentials() error {
	if c.Email == "" {
		return errors.New("email is required")
	}
	if c.Password == "" {
		return errors.New("password is required")
	}
	if !ValidateEmail(c.Email) {
		return errors.New("email is invalid")
	}
	return nil
}

type User struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"not null" json:"username" binding:"required"`
	Email     string    `gorm:"unique;not null" json:"email" binding:"required,email"`
	Password  string    `gorm:"not null" binding:"required" json:"password"`
	Books     []Book    `gorm:"foreignKey:OwnerID;references:ID;constraint:OnDelete:CASCADE" json:"books"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (u *User) ValidateUser() error {
	if u.Username == "" {
		return errors.New("username is required")
	}
	if u.Email == "" {
		return errors.New("email is required")
	}
	if u.Password == "" {
		return errors.New("password is required")
	}

	if !ValidateEmail(u.Email) {
		return errors.New("email is invalid")
	}

	return nil
}

type Book struct {
	ID         int       `gorm:"primaryKey" json:"id"`
	Title      string    `gorm:"not null" json:"title" binding:"required"`
	Edition    int       `json:"edition,omitempty"`
	Author     string    `gorm:"not null" json:"author" binding:"required"`
	PagesCount int       `gorm:"not null" json:"pages_count" binding:"required"`
	PagesRead  int       `gorm:"not null" json:"pages_read" binding:"required"`
	OwnerID    int       `json:"owner_id"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (b *Book) ValidateBook() error {
	if b.Title == "" {
		return errors.New("book title is required")
	}
	if b.Author == "" {
		return errors.New("book author is required")
	}
	if b.PagesCount <= 0 {
		return errors.New("invalid pages count")
	}
	if b.PagesRead < 0 || b.PagesRead > b.PagesCount {
		return errors.New("invalid pages read")
	}
	return nil
}
