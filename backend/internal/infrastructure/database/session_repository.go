package database

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/riii111/markdown-blog-api/internal/domain/model"
	"gorm.io/gorm"
)

type sessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) *sessionRepository {
	return &sessionRepository{db: db}
}

// セッションを作成
func (r *sessionRepository) Create(session *model.Session) error {
	if err := r.db.Create(session).Error; err != nil {
		log.Printf("Session creation failed: %v", err)
		return err
	}
	return nil
}

// IDによるセッション検索
func (r *sessionRepository) FindByID(id uuid.UUID) (*model.Session, error) {
	var session model.Session
	err := r.db.Preload("User").Where("id = ?", id).First(&session).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &session, nil
}

// トークンによるセッション検索
func (r *sessionRepository) FindByToken(ctx context.Context, token string) (*model.Session, error) {
	var session model.Session
	err := r.db.WithContext(ctx).Preload("User").Where("session_token = ?", token).First(&session).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &session, nil
}

// セッションを削除
func (r *sessionRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Session{}, "id = ?", id).Error
}

// 期限切れのセッションを削除
func (r *sessionRepository) DeleteExpired() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&model.Session{}).Error
}
