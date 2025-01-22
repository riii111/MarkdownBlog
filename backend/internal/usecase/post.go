package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/riii111/markdown-blog-api/internal/domain/model"
)

type PostUsecase struct {
	postRepo model.PostRepository
}

func NewPostUsecase(postRepo model.PostRepository) *PostUsecase {
	return &PostUsecase{
		postRepo: postRepo,
	}
}

// 投稿作成
func (u *PostUsecase) CreatePost(ctx context.Context, userID uuid.UUID, title, content string, seriesID *string) (*model.Post, error) {
	// スラッグの生成（タイトルから）
	slugStr := slug.Make(title)

	// UUIDの生成
	postID := uuid.New()

	var seriesUUID *uuid.UUID
	if seriesID != nil {
		parsed, err := uuid.Parse(*seriesID)
		if err != nil {
			return nil, fmt.Errorf("invalid series ID: %w", err)
		}
		seriesUUID = &parsed
	}

	post := &model.Post{
		BaseModel: model.BaseModel{
			ID: postID,
		},
		UserID:     userID,
		SeriesID:   seriesUUID,
		Title:      title,
		Content:    content,
		Slug:       fmt.Sprintf("%s-%s", slugStr, postID.String()[:8]),
		Status:     "draft",
		LikesCount: 0,
	}

	if err := u.postRepo.Create(post); err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	return post, nil
}

// 投稿削除
func (u *PostUsecase) DeletePost(ctx context.Context, userID, postID uuid.UUID) error {
	post, err := u.postRepo.FindByID(postID)
	if err != nil {
		return err
	}
	if post == nil {
		return errors.New("post not found")
	}

	// 投稿者のみが削除可能
	if post.UserID != userID {
		return errors.New("unauthorized")
	}

	return u.postRepo.Delete(postID)
}
