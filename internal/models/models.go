package models

import "time"

type User struct {
	ID	string 	`gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Login	string	`gorm:"type:varchar(255);uniqueIndex;not null"`
	Password	string 	`gorm:"column:password;type:text;not null"`
	CreatedAt	time.Time	`gorm:"not null;default:now()"`
}

func (User) TableName() string { return "users" }

type Token struct {
	Token	string	`gorm:"primaryKey;type:text"`
	UserID	string	`gorm:"column:user_id;type:uuid;not null;index"`
	CreatedAt time.Time	`gorm:"not null;default:now()"`
	ExpiresAt	time.Time	`gorm:"column:expires_at;not null;index"`
}

func (Token) TableName() string { return "tokens" }

type Document struct {
	ID	string	`gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID	string	`gorm:"column:user_id;type:uuid;not null;index"`
	Name	string	`gorm:"type:varchar(255);not null;index:idx_documents_name_created"`
	Mime	string	`gorm:"type:varchar(255);not null;default:''"`
	IsFile	bool	`gorm:"column:is_file;not null;default:false"`
	IsPublic	bool	`gorm:"column:is_public;not null;default:false"`
	FilePath	string	`gorm:"column:file_path;type:text"`
	JSONData	[]byte	`gorm:"column:json_data;type:jsonb"`
	CreatedAt	time.Time	`gorm:"not null;default:now();index:idx_documents_name_created"`

	Grant	[]string	`gorm:"-" json:"grant,omitempty"`
}

func (Document) TableName() string { return "documents" }

type DocumentGrant struct {
	DocumentID	string	`gorm:"column:document_id;type:uuid;primaryKey"`
	UserID	string	`gorm:"column:user_id;type:uuid;primaryKey;index"`
}

func (DocumentGrant) TableName() string { return "document_grants" }
