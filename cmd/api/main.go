package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"go-expense-tracker/db"
	"go-expense-tracker/handlers"
)

func main() {
	r := gin.Default()

	if err := db.Init(); err != nil {
		log.Fatal("couldnt connect to the database!")
	}

	r.POST("/login", handlers.Login)
	r.POST("/signup", handlers.Signup)

	authorized := r.Group("/")
	authorized.Use(handlers.CheckLogin())
	{
		authorized.POST("/signout", handlers.Signout)

		authorized.GET("/user", handlers.GetUsers)
		authorized.GET("/user/:id", handlers.GetUser)
		authorized.PUT("/user/:id", handlers.UpdateUser)
		authorized.DELETE("/user/:id", handlers.DeleteUser)

		authorized.GET("/account", handlers.GetAccounts)
		authorized.GET("/account/:id", handlers.GetAccount)
		authorized.POST("/account", handlers.CreateAccount)
		authorized.PUT("/account/:id", handlers.UpdateAccount)
		authorized.DELETE("/account/:id", handlers.DeleteAccount)

		authorized.GET("/transaction", handlers.GetTransactions)
		authorized.GET("/transaction/:id", handlers.GetTransaction)
		authorized.POST("/transaction", handlers.CreateTransaction)
		authorized.PUT("/transaction/:id", handlers.UpdateTransaction)
		authorized.DELETE("/transaction/:id", handlers.DeleteTransaction)

		authorized.GET("/transaction-category", handlers.GetTransactionCategories)
		authorized.GET("/transaction-category/:id", handlers.GetTransactionCategory)
		authorized.POST("/transaction-category", handlers.CreateTransactionCategory)
		authorized.PUT("/transaction-category/:id", handlers.UpdateTransactionCategory)
		authorized.DELETE("/transaction-category/:id", handlers.DeleteTransactionCategory)
	}

	r.Run(":2024")
}
