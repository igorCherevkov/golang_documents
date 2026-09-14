package store

import (
	"context"
	"documents/internal/models"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var filterColumns = map[string]string{
	"name":	"name",
	"mime":	"mime",
	"file":	"is_file",
	"public":	"is_public",
	"is_file":	"is_file",
	"is_public":	"is_public",
}

type DocumentStore struct {
	db *gorm.DB
}

func NewDocumentStore(db *gorm.DB) *DocumentStore {
	return &DocumentStore{db: db}
}

func (s *DocumentStore) Create(ctx context.Context, document models.Document, grantLogins []string) (models.Document, error) {
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&document).Error; err != nil {
			return fmt.Errorf("[store] error while insert document: %w", err)
		}

		if len(grantLogins) == 0 {
			return nil
		}

		var userIDs []string
		if err := tx.Model(&models.User{}).Where("login IN ?", grantLogins).Pluck("id", &userIDs).Error; err != nil {
			return fmt.Errorf("[store] error while resolve grant login")
		}

		if len(userIDs) == 0 {
			return nil
		}

		grants := make([]models.DocumentGrant, 0, len(userIDs))
		
		for _, userID := range userIDs {
			grants = append(grants, models.DocumentGrant{DocumentID: document.ID, UserID: userID})
		}

		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&grants).Error; err != nil {
			return fmt.Errorf("[store] error while insert grants: %w", err)
		}

		return nil
	})

	if err != nil {
		return models.Document{}, err
	}

	document.Grant = grantLogins

	return document, nil
}

type ListParams struct {
	RequesterUserID	string
	UserLogin	string
	FilterKey	string
	FilterValue	string
	Limit	int
}
 
func (s *DocumentStore) List(ctx context.Context, params ListParams) ([]models.Document, error) {
	query := s.db.WithContext(ctx).Model(&models.Document{})

	if params.UserLogin != "" {
		query = query.Joins("JOIN users owner ON owner.id = documents.user_id").
			Where("owner.login = ?", params.UserLogin).
			Where(
				"(owner.id = ? OR documents.is_public = true OR EXISTS (" +
				"SELECT 1 FROM document_grants g WHERE g.document_id = documents.id AND g.user_id = ?" + 
				"))",
				params.RequesterUserID, params.RequesterUserID,
			)
	} else {
		query = query.Where("documents.user_id = ?", params.RequesterUserID)
	}

	if params.FilterKey != "" {
		column, ok := filterColumns[params.FilterKey]
		if !ok {
			return nil, ErrBadFilterKey
		}

		query.Where(fmt.Sprintf("documents.%s::text = ?", column), params.FilterValue)
	}

	query = query.Order("documents.name ASC, documents.created_at ASC")
	if params.Limit > 0 {
		query = query.Limit(params.Limit)
	}

	var docs []models.Document
	if err := query.Find(&docs).Error; err != nil {
		return nil, fmt.Errorf("store: list documents: %w", err)
	}

	if len(docs) == 0 {
		return docs, nil
	}

	ids := make([]string, len(docs))
	for i, doc := range docs {
		ids[i] = doc.ID
	}

	grantsByDoc, err := s.grantLoginsByDocID(ctx, ids)
	if err != nil {
		return nil, err
	}

	for i := range docs {
		docs[i].Grant = grantsByDoc[docs[i].ID]
	}

	return docs, nil
}

func (s *DocumentStore) GetByID(ctx context.Context, id string) (models.Document, error) {
	var document models.Document

	err := s.db.WithContext(ctx).Where("id = ?", id).First(&document).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Document{}, ErrNotFound
		}

		return models.Document{}, fmt.Errorf("[store] error while get document: %w", err)
	}

	grantsByDoc, err := s.grantLoginsByDocID(ctx, []string{id})
	if err != nil {
		return models.Document{}, err
	}

	document.Grant = grantsByDoc[id]

	return document, nil
}

func (s *DocumentStore) Delete(ctx context.Context, id string) error {
	result := s.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Document{})
	if result.Error != nil {
		return fmt.Errorf("[store] error while delete document: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *DocumentStore) grantLoginsByDocID(ctx context.Context, docIDs []string) (map[string][]string, error) {
	type row struct {
		DocumentID string
		Login      string
	}

	var rows []row

	err := s.db.WithContext(ctx).
		Table("document_grants AS g").
		Select("g.document_id AS document_id, u.login AS login").
		Joins("JOIN users u ON u.id = g.user_id").
		Where("g.document_id IN ?", docIDs).
		Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("[store] error while load grants: %w", err)
	}

	result := make(map[string][]string, len(docIDs))

	for _, row := range rows {
		result[row.DocumentID] = append(result[row.DocumentID], row.Login)
	}

	return result, nil
}