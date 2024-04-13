package models

import (
	"time"

	"gorm.io/gorm"
)

type TransactionTypes string

const (
	Expense TransactionTypes = "expense"
	Income  TransactionTypes = "income"
)

type Transaction struct {
	gorm.Model
	AccountID   uint
	CategoryID  uint
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Value       int64            `json:"value"`
	Date        time.Time        `json:"date"`
	Type        TransactionTypes `json:"type"        sql:"type:ENUM('expense','income')"`
	Completed   bool             `json:"completed"`
}
