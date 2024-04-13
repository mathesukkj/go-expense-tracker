package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"go-expense-tracker/db"
)

func main() {
	r := gin.Default()

	if err := db.Init(); err != nil {
		log.Fatal("couldnt connect to the database!")
	}

	r.Run(":2024")
}
