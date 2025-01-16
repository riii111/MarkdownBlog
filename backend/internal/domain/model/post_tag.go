package model

import (
	"github.com/google/uuid"
	"time"
)

// PostTag 記事とタグの中間テーブル
type PostTag struct {
	PostID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	TagID     uuid.UUID `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time `gorm:"not null"`
	Post      Post      `gorm:"foreignKey:PostID"`
	Tag       Tag       `gorm:"foreignKey:TagID"`
}
