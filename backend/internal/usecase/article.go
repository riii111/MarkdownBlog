package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/riii111/markdown-blog-api/internal/domain/model"
)

// ページネーションのデフォルト値
const (
	DefaultPage    = 1
	DefaultPerPage = 9
	MaxPerPage     = 100
)

type ArticleUsecase struct {
	articleRepo model.ArticleRepository
}

func NewArticleUsecase(articleRepo model.ArticleRepository) *ArticleUsecase {
	return &ArticleUsecase{
		articleRepo: articleRepo,
	}
}

// 記事作成
func (u *ArticleUsecase) CreateArticle(ctx context.Context, userID uuid.UUID, title, content string, seriesID *string) (*model.Article, error) {
	articleID := uuid.New()

	// UUIDからハイフンを除去して最初の13文字を使用
	slugStr := strings.ReplaceAll(articleID.String(), "-", "")[:13]

	var seriesUUID *uuid.UUID
	if seriesID != nil {
		parsed, err := uuid.Parse(*seriesID)
		if err != nil {
			return nil, fmt.Errorf("invalid series ID: %w", err)
		}
		seriesUUID = &parsed
	}

	article := &model.Article{
		BaseModel: model.BaseModel{
			ID: articleID,
		},
		UserID:     userID,
		SeriesID:   seriesUUID,
		Title:      title,
		Content:    content,
		Slug:       slugStr,
		Status:     "draft",
		LikesCount: 0,
	}

	if err := u.articleRepo.Create(article); err != nil {
		return nil, fmt.Errorf("failed to create article: %w", err)
	}

	return article, nil
}

// Slug使って記事削除
func (u *ArticleUsecase) DeleteArticleBySlug(ctx context.Context, userID uuid.UUID, slug string) error {
	article, err := u.articleRepo.FindBySlug(slug)
	if err != nil {
		return err
	}
	if article == nil {
		return errors.New("article not found")
	}

	if article.UserID != userID {
		return errors.New("unauthorized")
	}

	return u.articleRepo.Delete(article.ID)
}

// GetPublishedArticles 公開済み記事一覧を取得
func (u *ArticleUsecase) GetPublishedArticles(ctx context.Context, page, perPage int) ([]model.Article, int, error) {
	return u.articleRepo.FindPublished(page, perPage)
}

// GetArticleBySlug slugで記事の詳細を取得
func (u *ArticleUsecase) GetArticleBySlug(ctx context.Context, slug string) (*model.Article, error) {
	article, err := u.articleRepo.FindBySlugWithRelations(slug)
	if err != nil {
		return nil, fmt.Errorf("failed to find article: %w", err)
	}
	if article == nil {
		return nil, errors.New("article not found")
	}

	return article, nil
}
