package model

import (
	"time"

	"github.com/google/uuid"
)

// Session セッションモデル
type Session struct {
	BaseModel
	UserID       uuid.UUID `gorm:"type:uuid;not null"`
	SessionToken string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	ExpiresAt    time.Time `gorm:"not null"`
	User         User      `gorm:"foreignKey:UserID"`
}

// SessionRepository セッションリポジトリのインターフェース
type SessionRepository interface {
	Create(session *Session) error
	FindByID(id uuid.UUID) (*Session, error)
	FindByToken(token string) (*Session, error)
	Delete(id uuid.UUID) error
	DeleteExpired() error
}
