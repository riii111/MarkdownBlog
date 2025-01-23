package database

import (
	"fmt"
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

// slugで記事を検索し、関連データも取得
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

// カーソルベースで公開済み記事一覧を取得
func (r *articleRepository) FindPublished(limit int, cursor *string) ([]model.Article, *string, error) {
	var articles []model.Article
	query := r.db.Preload("User").
		Where("status = ? AND published_at <= ?", "published", time.Now())

	// カーソルが指定されている場合、そのIDより小さいものを取得
	if cursor != nil {
		cursorUUID, err := uuid.Parse(*cursor)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid cursor format: %w", err)
		}
		query = query.Where("id < ?", cursorUUID)
	}

	// ID降順で取得（新しい順）
	result := query.Order("id DESC").
		Limit(limit + 1). // 次ページの有無を確認するため、1つ多めに取得
		Find(&articles)

	if result.Error != nil {
		return nil, nil, result.Error
	}

	var nextCursor *string
	hasMore := len(articles) > limit

	// 次ページがある場合
	if hasMore {
		// 最後の要素を除去（limit + 1件目）
		articles = articles[:limit]
		// 次のカーソルを設定
		lastID := articles[len(articles)-1].ID.String()
		nextCursor = &lastID
	}

	return articles, nextCursor, nil
}
