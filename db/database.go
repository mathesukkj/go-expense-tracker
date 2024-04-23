package db

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"go-expense-tracker/models"
)

var Gorm *gorm.DB

func Init() error {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=America/Sao_Paulo",
		dbHost,
		dbPort,
		dbUser,
		dbPassword,
		dbName,
	)
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
