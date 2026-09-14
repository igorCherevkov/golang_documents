package store

import (
	"context"
	"documents/internal/models"
	"errors"
	"time"

	"gorm.io/gorm"
)

type TokenStore struct {
	db *gorm.DB
}

func NewTokenStore(db *gorm.DB) *TokenStore {
	return &TokenStore{db: db}
}

func (s *TokenStore) Create(ctx context.Context, token, userID string, ttl time.Duration) error {
	t := models.Token{
		Token: token,
		UserID: userID,
		ExpiresAt: time.Now().Add(ttl),
	}

	return s.db.WithContext(ctx).Create(&t).Error
}

func (s *TokenStore) UserIDByToken(ctx context.Context, token string) (string, error) {
	var t models.Token

	err := s.db.WithContext(ctx).Where("token = ? AND expires_at > ?", token, time.Now()).First(&t).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrNotFound
		}

		return "", err
	}

	return t.UserID, nil
}

func (s *TokenStore) DeleteSession(ctx context.Context, token string) error {
	result := s.db.WithContext(ctx).Where("token = ?", token).Delete(&models.Token{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *TokenStore) DeleteExpiredTokens(ctx context.Context) (int64, error) {
	result := s.db.WithContext(ctx).Where("expires_at <= ?", time.Now()).Delete(&models.Token{})

	return result.RowsAffected, result.Error
}