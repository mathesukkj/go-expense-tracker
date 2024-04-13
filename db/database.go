package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"go-expense-tracker/models"
)

var Gorm *gorm.DB

func Init() error {
	dsn := "host=localhost user=postgres password=password dbname=expense_tracker port=5432 sslmode=disable TimeZone=America/Sao_Paulo"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	db.AutoMigrate(&models.Account{})
	db.AutoMigrate(&models.Transaction{})
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.TransactionCategory{})

	Gorm = db

	return nil
}
