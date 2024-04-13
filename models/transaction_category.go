package models

import "gorm.io/gorm"

type TransactionCategory struct {
	gorm.Model
	IconUrl      string `json:"iconUrl"`
	ColorHex     string `json:"colorHex"`
	Name         string `json:"name"`
	Transactions []Transaction
}
