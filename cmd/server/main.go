package main

import (
	"context"
	"documents/internal/api"
	"documents/internal/config"
	"documents/internal/db"
	"documents/internal/store"
	"documents/internal/token_cleaning"
	"log"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.Load()

	connection, err := db.New(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("[server] error while connect db: %v", err)
	}

	users := store.NewUserStore(connection)
	tokens := store.NewTokenStore(connection)
	documents := store.NewDocumentStore(connection)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	token_cleaning.StartTokenCleaning(ctx, tokens, 1*time.Hour)

	app := api.NewApp(cfg, users, tokens, documents)

	go func() {
		<-ctx.Done()
		log.Println("[server] shutting down gracefully")
		if err := app.ShutdownWithContext(context.Background()); err != nil {
			log.Printf("[server] shutdown error: %v", err)
		}
	}()

	log.Printf("[server] listening on :%s", cfg.Port)

	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Printf("[server] listen error: %v", err)
	}

	log.Println("[server] stopped")
}
