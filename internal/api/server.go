package api

import (
	"documents/internal/cache"
	"documents/internal/config"
	"documents/internal/store"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type Server struct {
	cfg       config.Config
	users     *store.UserStore
	tokens    *store.TokenStore
	documents *store.DocumentStore
	validate  *validator.Validate
	cache	*cache.Cache
}

func NewServer(cfg config.Config, users *store.UserStore, tokens *store.TokenStore, documents *store.DocumentStore) *Server {
	return &Server{
		cfg:       cfg,
		users:     users,
		tokens:    tokens,
		documents: documents,
		validate:  validator.New(),
		cache: cache.NewCache(),
	}
}

func (s *Server) RegisterRoutes(app *fiber.App) {
	auth := requireAuth(s.tokens)

	app.Post("/api/register", s.register)
	app.Post("/api/auth", s.auth)
	app.Delete("/api/auth/logout", auth, s.logout)

	app.Post("/api/docs", auth, s.uploadDocument)
	app.Get("/api/docs", auth, s.listDocuments)
	app.Head("/api/docs", auth, s.listDocuments)
	app.Get("/api/docs/:id", auth, s.getDocument)
	app.Head("/api/docs/:id", auth, s.getDocument)
	app.Delete("/api/docs/:id", auth, s.deleteDocument)
}
