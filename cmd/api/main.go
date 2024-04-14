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

	r.GET("/user", handlers.GetUsers)
	r.GET("/user/:id", handlers.GetUser)
	r.PUT("/user/:id", handlers.UpdateUser)
	r.DELETE("/user/:id", handlers.DeleteUser)

	r.POST("/login", handlers.Login)
	r.POST("/signup", handlers.Signup)

	r.Run(":2024")
}
