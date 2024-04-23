package main

import (
	"flag"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"

	"go-expense-tracker/db"
	"go-expense-tracker/handlers"
	cache "go-expense-tracker/redis"
)

func main() {
	r := gin.Default()

	if err := db.Init(); err != nil {
		log.Fatal("couldnt connect to the database!")
	}

	cache.Init()

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

		authorized.GET("/dashboard/current-balance", handlers.GetCurrentBalance)
		authorized.GET("/dashboard/total-income", handlers.GetTotalIncome)
		authorized.GET("/dashboard/total-expense", handlers.GetTotalExpense)
	}

	port := flag.String("port", "8080", "define the port of the app")
	flag.Parse()

	r.Run(":" + *port)
}
