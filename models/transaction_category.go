package models

import "gorm.io/gorm"

type TransactionCategory struct {
	gorm.Model
	UserID       uint          `json:"-"`
	IconUrl      string        `json:"icon_url"`
	ColorHex     string        `json:"color_hex"`
	Name         string        `json:"name"`
	Transactions []Transaction `json:"-"`
}

type TransactionCategoryPayload struct {
	IconUrl  string `json:"icon_url"  binding:"required"`
	ColorHex string `json:"color_hex" binding:"required"`
	Name     string `json:"name"      binding:"required"`
}
