package database

import (
	"time"

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

// slugによる投稿検索
func (r *articleRepository) FindBySlug(slug string) (*model.Article, error) {
	var article model.Article
	result := r.db.Where("slug = ?", slug).First(&article)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &article, nil
}

// FindPublished 公開済み記事一覧を取得
func (r *articleRepository) FindPublished(page, perPage int) ([]model.Article, int, error) {
	var articles []model.Article
	var total int64

	// 総件数を取得
	if err := r.db.Model(&model.Article{}).
		Where("status = ? AND published_at <= ?", "published", time.Now()).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 記事一覧を取得
	result := r.db.Preload("User").
		Where("status = ? AND published_at <= ?", "published", time.Now()).
		Order("published_at DESC").
		Offset((page - 1) * perPage).
		Limit(perPage).
		Find(&articles)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	return articles, int(total), nil
}

// FindBySlugWithRelations slugで記事を検索し、関連データも取得
func (r *articleRepository) FindBySlugWithRelations(slug string) (*model.Article, error) {
	var article model.Article
	result := r.db.Preload("User").
		Preload("Series").
		Preload("Tags").
		Where("slug = ?", slug).
		First(&article)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}

	return &article, nil
}
