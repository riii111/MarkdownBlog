package model

import "github.com/google/uuid"

// User ユーザーモデル
type User struct {
	BaseModel
	Email        string `gorm:"type:varchar(255);uniqueIndex;not null"`
	PasswordHash string `gorm:"type:varchar(255);not null"`
	DisplayName  string `gorm:"type:varchar(100);not null"`
}

// UserRepository ユーザーリポジトリのインターフェース
type UserRepository interface {
	Create(user *User) error
	FindByID(id uuid.UUID) (*User, error)
	FindByEmail(email string) (*User, error)
	Update(user *User) error
	Delete(id uuid.UUID) error
}
