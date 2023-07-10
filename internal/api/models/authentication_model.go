package models

type UserCredentials struct {
	Email    string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}
