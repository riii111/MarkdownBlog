package model

import (
	"time"

	"github.com/google/uuid"
)

// 記事モデル
type Article struct {
	BaseModel
	UserID      uuid.UUID  `gorm:"type:uuid;not null"`
	SeriesID    *uuid.UUID `gorm:"type:uuid"`
	Title       string     `gorm:"type:varchar(255);not null"`
	Content     string     `gorm:"type:text;not null"`
	Slug        string     `gorm:"type:varchar(255);uniqueIndex;not null"`
	Status      string     `gorm:"type:varchar(20);not null;default:'draft';check:status IN ('draft','published')"`
	LikesCount  int        `gorm:"not null;default:0"`
	PublishedAt *time.Time
	User        User    `gorm:"foreignKey:UserID"`
	Series      *Series `gorm:"foreignKey:SeriesID"`
	Tags        []Tag   `gorm:"many2many:article_tags;"`
}

// 記事リポジトリのインターフェース
type ArticleRepository interface {
	Create(article *Article) error
	FindByID(id uuid.UUID) (*Article, error)
	FindBySlug(slug string) (*Article, error)
	Delete(id uuid.UUID) error
	FindPublished(limit int, cursor *string) ([]Article, *string, error)
	FindBySlugWithRelations(slug string) (*Article, error)
	FindByUserID(userID uuid.UUID, limit int, cursor *string, status string) ([]Article, *string, error)
	FindByTagSlug(tagSlug string, limit int, cursor *string) ([]Article, *string, error)
}
