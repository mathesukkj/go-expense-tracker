package models

import "gorm.io/gorm"

type Account struct {
	gorm.Model
	UserID       uint
	Name         string `json:"name"`
	Balance      int64  `json:"balance"`
	Transactions []Transaction
}
