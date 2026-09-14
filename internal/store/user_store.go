package store

import (
	"context"
	"documents/internal/models"
	"errors"

	"gorm.io/gorm"
)

type UserStore struct {
	db *gorm.DB
}

func NewUserStore(db *gorm.DB) *UserStore {
	return &UserStore{db: db}
}

func (s *UserStore) Create(ctx context.Context, login, password string) (models.User, error) {
	u := models.User{Login: login, Password: password}

	err := s.db.WithContext(ctx).Create(&u).Error
	if err != nil {
		if isUniqueViolation(err) {
			return models.User{}, ErrLoginTaken
		}

		return models.User{}, err
	}

	return u, nil
}

func (s *UserStore) GetUserByLogin(ctx context.Context, login string) (models.User, error) {
	var u models.User

	err := s.db.WithContext(ctx).Where("login = ?", login).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.User{}, ErrNotFound
		}

		return models.User{}, err
	}

	return u, nil
}

func (s *UserStore) GetUserById(ctx context.Context, id string) (models.User, error) {
	var u models.User

	err := s.db.WithContext(ctx).Where("id = ?", id).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.User{}, ErrNotFound
		}

		return models.User{}, err
	}

	return u, nil
}
