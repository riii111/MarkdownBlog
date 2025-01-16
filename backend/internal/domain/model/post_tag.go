package model

import "github.com/google/uuid"

// PostTag 記事とタグの中間テーブル
type PostTag struct {
	BaseModel
	PostID uuid.UUID `gorm:"type:uuid;primaryKey"`
	TagID  uuid.UUID `gorm:"type:uuid;primaryKey"`
	Post   Post      `gorm:"foreignKey:PostID"`
	Tag    Tag       `gorm:"foreignKey:TagID"`
}
