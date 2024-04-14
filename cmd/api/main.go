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
		authorized.GET("/user", handlers.GetUsers)
		authorized.GET("/user/:id", handlers.GetUser)
		authorized.PUT("/user/:id", handlers.UpdateUser)
		authorized.DELETE("/user/:id", handlers.DeleteUser)
	}

	r.Run(":2024")
}
