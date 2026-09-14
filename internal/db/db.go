package db

import (
	"documents/internal/models"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func New(databaseURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("[db] error while open connection: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("[db] error while get sql.db: %w", err)
	}
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("[db] error while ping: %w", err)
	}

	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS pgcrypto`).Error; err != nil {
		return fmt.Errorf("[db] error while create extension pgcrypto: %w", err)
	}

	return db.AutoMigrate(
		&models.User{},
		&models.Token{},
		&models.Document{},
		&models.DocumentGrant{},
	)
}
