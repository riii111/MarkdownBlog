package database

import (
	"github.com/google/uuid"
	"github.com/riii111/markdown-blog-api/internal/domain/model"
	"gorm.io/gorm"
)

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *postRepository {
	return &postRepository{db: db}
}

// 投稿作成
func (r *postRepository) Create(post *model.Post) error {
	result := r.db.Create(post)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// IDによる投稿検索
func (r *postRepository) FindByID(id uuid.UUID) (*model.Post, error) {
	var post model.Post
	result := r.db.Where("id = ?", id).First(&post)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &post, nil
}

// 投稿削除
func (r *postRepository) Delete(id uuid.UUID) error {
	result := r.db.Delete(&model.Post{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
