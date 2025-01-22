package model

import "github.com/google/uuid"

// Series シリーズモデル
type Series struct {
	BaseModel
	UserID      uuid.UUID `gorm:"type:uuid;not null"`
	Title       string    `gorm:"type:varchar(255);not null"`
	Description *string   `gorm:"type:text"`
	Slug        string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	User        User      `gorm:"foreignKey:UserID"`
	Articles    []Article `gorm:"foreignKey:SeriesID"`
}

// SeriesRepository シリーズリポジトリのインターフェース
type SeriesRepository interface {
	Create(series *Series) error
	FindByID(id uuid.UUID) (*Series, error)
	FindBySlug(slug string) (*Series, error)
	Update(series *Series) error
	Delete(id uuid.UUID) error
}
