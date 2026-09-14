package main

import (
	"documents/internal/config"
	"documents/internal/db"
	"log"
)

func main() {
	cfg := config.Load()

	connection, err := db.New(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("[migrate] error while connect: %v", err)
	}

	if err := db.AutoMigrate(connection); err != nil {
		log.Fatalf("[migrate] error while automigrate: %v", err)
	}

	log.Println("[migrate] done")
}
