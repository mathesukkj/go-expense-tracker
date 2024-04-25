package main

import (
	"flag"
	"log"

	_ "github.com/joho/godotenv/autoload"

	"go-expense-tracker/internal/db"
	cache "go-expense-tracker/internal/redis"
	"go-expense-tracker/internal/routes"
)

func main() {
	r := routes.NewRouter()

	if err := db.Init(); err != nil {
		log.Fatal("couldnt connect to the database!")
	}

	cache.Init()

	port := flag.String("port", "8080", "define the port of the app")
	flag.Parse()

	r.Run(":" + *port)
}
