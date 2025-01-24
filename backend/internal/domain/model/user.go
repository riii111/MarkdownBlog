package model

import (
	"context"

	"github.com/google/uuid"
)

// ユーザーモデル
type User struct {
	BaseModel
	Email        string `gorm:"type:varchar(255);uniqueIndex;not null"`
	PasswordHash string `gorm:"type:varchar(255);not null"`
	DisplayName  string `gorm:"type:varchar(100);not null"`
}

// ユーザーリポジトリのインターフェース
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uuid.UUID) error
}
