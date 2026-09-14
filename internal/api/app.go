package api

import (
	"documents/internal/config"
	"documents/internal/store"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func NewApp(cfg config.Config, users *store.UserStore, tokens *store.TokenStore, documents *store.DocumentStore) *fiber.App {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	app.Use(requestLogger())

	app.Get("/health", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"success":   true,
			"message":   "Documents API is running",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})

	server := NewServer(cfg, users, tokens, documents)
	server.RegisterRoutes(app)

	return app
}

func isMultipartRequest(ctx *fiber.Ctx) bool {
	return strings.HasPrefix(ctx.Get(fiber.HeaderContentType), "multipart/form-data")
}

func requestLogger() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		start := time.Now()

		err := ctx.Next()

		status := ctx.Response().StatusCode()
		duration := time.Since(start)
		path := ctx.Path()

		switch {
		case isMultipartRequest(ctx):
			log.Printf("%s %s -> %d (%s) ip=%s body=<multipart, skipped>", ctx.Method(), path, status, duration, ctx.IP())
		default:
			log.Printf("%s %s -> %d (%s) ip=%s body=%s", ctx.Method(), path, status, duration, ctx.IP(), ctx.Body())
		}

		return err
	}
}
