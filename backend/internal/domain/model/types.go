package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// 共通のフィールドとメソッドを持つ基底モデル
type BaseModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

// UUIDv7を生成する共通処理
func (b *BaseModel) BeforeCreate(tx *gorm.DB) error {
	// IDが未設定の場合のみ生成
	if b.ID == uuid.Nil {
		id, err := uuid.NewV7()
		if err != nil {
			return err
		}
		b.ID = id
	}
	return nil
}
