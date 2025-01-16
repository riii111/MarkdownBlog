package model

import "github.com/google/uuid"

// Tag タグモデル
type Tag struct {
	BaseModel
	Name  string `gorm:"type:varchar(50);uniqueIndex;not null"`
	Slug  string `gorm:"type:varchar(50);uniqueIndex;not null"`
	Posts []Post `gorm:"many2many:post_tags;"`
}

// TagRepository タグリポジトリのインターフェース
type TagRepository interface {
	Create(tag *Tag) error
	FindByID(id uuid.UUID) (*Tag, error)
	FindBySlug(slug string) (*Tag, error)
	FindAll() ([]Tag, error)
	Update(tag *Tag) error
	Delete(id uuid.UUID) error
}
