package model

import (
	"time"

	"github.com/google/uuid"
)

// Post 記事モデル
type Post struct {
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
	Tags        []Tag   `gorm:"many2many:post_tags;"`
}

// PostRepository 記事リポジトリのインターフェース
type PostRepository interface {
	Create(post *Post) error
	FindByID(id uuid.UUID) (*Post, error)
	Delete(id uuid.UUID) error
}
