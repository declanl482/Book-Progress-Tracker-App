package models

type UserCredentials struct {
	Email    string `form:"username" binding:"required,email"`
	Password string `form:"password" binding:"required"`
}
