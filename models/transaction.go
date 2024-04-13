package models

import (
	"time"

	"gorm.io/gorm"
)

type transactionTypes string

const (
	Expense transactionTypes = "expense"
	Income  transactionTypes = "income"
)

type Transaction struct {
	gorm.Model
	AccountID             uint
	TransactionCategoryID uint
	Name                  string           `json:"name"`
	Description           string           `json:"description"`
	Value                 int64            `json:"value"`
	Date                  time.Time        `json:"date"`
	Type                  transactionTypes `json:"type"        sql:"type:ENUM('expense','income')"`
	Completed             bool             `json:"completed"`
}
