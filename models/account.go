package models

import "gorm.io/gorm"

type Account struct {
	gorm.Model
	UserID       uint          `json:"user_id" gorm:"index"`
	Name         string        `json:"name"`
	Balance      int64         `json:"balance"`
	Transactions []Transaction `json:"-"`
}

type AccountPayload struct {
	UserID  uint   `json:"-"`
	Name    string `json:"name"    binding:"required"`
	Balance int64  `json:"balance"                    default:"0"`
}
