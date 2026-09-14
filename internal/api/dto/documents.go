package dto

import "documents/internal/models"

type DocMeta struct {
	Name   string   `json:"name" validate:"required"`
	File   bool     `json:"file"`
	Public bool     `json:"public"`
	Token  string   `json:"token"`
	Mime   string   `json:"mime" validate:"required"`
	Grant  []string `json:"grant" validate:"omitempty,dive,alphanum"`
}

type DocResponse struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Mime    string   `json:"mime"`
	File    bool     `json:"file"`
	Public  bool     `json:"public"`
	Created string   `json:"created"`
	Grant   []string `json:"grant,omitempty"`
}

type ListDocsQuery struct {
	Login string `query:"login"`
	Key   string `query:"key" validate:"omitempty,oneof=name mime file public"`
	Value string `query:"value"`
	Limit int    `query:"limit" validate:"omitempty,min=0,max=1000"`
}

func ToDocResponse(documents models.Document) DocResponse {
	return DocResponse{
		ID:      documents.ID,
		Name:    documents.Name,
		Mime:    documents.Mime,
		File:    documents.IsFile,
		Public:  documents.IsPublic,
		Created: documents.CreatedAt.Format("2006-01-02 15:04:05"),
		Grant:   documents.Grant,
	}
}
