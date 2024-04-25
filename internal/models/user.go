package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string    `json:"name"`
	Email    string    `json:"email" gorm:"index"`
	Password string    `json:"-"`
	Accounts []Account `json:"-"`
}

type UserPayload struct {
	Name     string `json:"name"     binding:"required"`
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
