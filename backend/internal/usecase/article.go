package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/riii111/markdown-blog-api/internal/domain/model"
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
	slugStr := slug.Make(title)
	if slugStr == "" {
		slugStr = "untitled" // タイトルが空の場合のデフォルト値
	}

	articleID := uuid.New()

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
		Slug:       fmt.Sprintf("%s-%s", slugStr, articleID.String()[:8]),
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
