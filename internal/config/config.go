package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string
	AdminToken  string
	StorageDir  string
	TokenTTL	time.Duration
}

func Load() Config {
	if err := godotenv.Load(); err != nil {
		log.Println("[config]: .env не найден")
	}

	return Config{
		Port:        getEnv("PORT", "8000"),
		DatabaseURL: getEnv("DATABASE_URL", "postgresql://postgres:postgres@localhost:5432/documents?schema=public"),
		AdminToken:  getEnv("ADMIN_TOKEN", "secret_admin_token"),
		StorageDir:  getEnv("STORAGE_DIR", "./storage"),
		TokenTTL: getEnvDuration("TOKEN_TTL", 24 * time.Hour),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		log.Printf("config: invalid %s=%q, using default %s", key, v, fallback)
		return fallback
	}
	return d
}
