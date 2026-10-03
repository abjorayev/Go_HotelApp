package main

import (
	"fmt"
	"hotel_app/config"
	"hotel_app/db"
	"log"
)

func main() {
	cfg := config.Load()

	database, err := db.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Migrate(database)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database connected and migrated!")
}