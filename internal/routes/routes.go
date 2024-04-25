package routes

import (
	"github.com/gin-gonic/gin"

	"go-expense-tracker/internal/handlers"
	"go-expense-tracker/internal/middleware"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/login", handlers.Login)
	r.POST("/signup", handlers.Signup)

	authorized := r.Group("/")
	authorized.Use(middleware.CheckLogin())
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
		authorized.GET("/dashboard/month-balance", handlers.GetMonthlyBalance)
	}

	return r
}
