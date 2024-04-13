package main

import (
	"log"

	"gotodo/db"
	"gotodo/handlers"
)

func main() {
	if err := db.Init(); err != nil {
		log.Fatal("couldnt connect to the database!")
	}
}
