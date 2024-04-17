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
	AccountID             uint             `json:"account_id"`
	TransactionCategoryID uint             `json:"category_id"`
	Description           string           `json:"description"`
	Value                 int64            `json:"value"`
	Date                  time.Time        `json:"date"`
	Type                  transactionTypes `json:"type"        sql:"type:ENUM('expense','income')"`
	Completed             bool             `json:"completed"`
}

type TransactionPayload struct {
	AccountID             uint             `json:"account_id"  binding:"required"`
	TransactionCategoryID uint             `json:"category_id" binding:"required"`
	Description           string           `json:"description" binding:"required"`
	Value                 int64            `json:"value"       binding:"required"`
	Date                  time.Time        `json:"date"        binding:"required"`
	Type                  transactionTypes `json:"type"        binding:"required,oneof=expense income"`
	Completed             bool             `json:"completed"                                           default:"false"`
}
