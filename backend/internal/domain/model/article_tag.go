package model

import "github.com/google/uuid"

// 記事とタグの中間テーブル
type ArticleTag struct {
	BaseModel
	ArticleID uuid.UUID `gorm:"type:uuid;primaryKey"`
	TagID     uuid.UUID `gorm:"type:uuid;primaryKey"`
	Article   Article   `gorm:"foreignKey:ArticleID"`
	Tag       Tag       `gorm:"foreignKey:TagID"`
}
