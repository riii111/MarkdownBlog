package database

import (
	"github.com/google/uuid"
	"github.com/riii111/markdown-blog-api/internal/domain/model"
	"gorm.io/gorm"
)

type articleRepository struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) *articleRepository {
	return &articleRepository{db: db}
}

// 投稿作成
func (r *articleRepository) Create(article *model.Article) error {
	result := r.db.Create(article)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// IDによる投稿検索
func (r *articleRepository) FindByID(id uuid.UUID) (*model.Article, error) {
	var article model.Article
	result := r.db.Where("id = ?", id).First(&article)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &article, nil
}

// 投稿削除
func (r *articleRepository) Delete(id uuid.UUID) error {
	result := r.db.Delete(&model.Article{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
