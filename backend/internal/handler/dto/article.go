package dto

import (
	"time"

	"github.com/riii111/markdown-blog-api/internal/domain/model"
)

type CreateArticleRequest struct {
	// 空のリクエストボディを表現するため、フィールドなし
}

type CreateArticleResponse struct {
	Slug string `json:"slug" example:"slug"`
}

// カーソルベースページネーション用のメタ情報
type CursorPaginationMeta struct {
	NextCursor   *string `json:"next_cursor,omitempty"`
	HasMore      bool    `json:"has_more"`
	ItemsPerPage int     `json:"items_per_page"`
}

// 記事一覧のレスポンス
type ArticleListResponse struct {
	Data       []ArticleListItem    `json:"data"`
	Pagination CursorPaginationMeta `json:"pagination"`
}

// 記事一覧のアイテム
type ArticleListItem struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Slug        string    `json:"slug"`
	User        UserInfo  `json:"user"`
	LikesCount  int       `json:"likes_count"`
	PublishedAt time.Time `json:"published_at"`
}

// 記事詳細のレスポンス
type ArticleDetailResponse struct {
	Data ArticleDetail `json:"data"`
}

type ArticleDetail struct {
	ID          string      `json:"id"`
	Title       string      `json:"title"`
	Content     string      `json:"content"`
	Slug        string      `json:"slug"`
	Status      string      `json:"status"`
	User        UserInfo    `json:"user"`
	Series      *SeriesInfo `json:"series,omitempty"`
	Tags        []TagInfo   `json:"tags"`
	LikesCount  int         `json:"likes_count"`
	PublishedAt *time.Time  `json:"published_at,omitempty"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// 記事一覧のレスポンスDTO変換
func NewArticleListResponse(articles []model.Article, limit int, nextCursor *string) ArticleListResponse {
	items := make([]ArticleListItem, len(articles))
	for i, article := range articles {
		items[i] = NewArticleListItem(article)
	}

	return ArticleListResponse{
		Data: items,
		Pagination: CursorPaginationMeta{
			NextCursor:   nextCursor,
			HasMore:      nextCursor != nil,
			ItemsPerPage: limit,
		},
	}
}

// 記事一覧のアイテムDTO変換
func NewArticleListItem(article model.Article) ArticleListItem {
	return ArticleListItem{
		ID:          article.ID.String(),
		Title:       article.Title,
		Slug:        article.Slug,
		LikesCount:  article.LikesCount,
		PublishedAt: *article.PublishedAt,
		User: UserInfo{
			ID:          article.User.ID.String(),
			DisplayName: article.User.DisplayName,
		},
	}
}

// 記事詳細のDTO変換
func NewArticleDetailResponse(article *model.Article) ArticleDetailResponse {
	detail := ArticleDetail{
		ID:          article.ID.String(),
		Title:       article.Title,
		Content:     article.Content,
		Slug:        article.Slug,
		Status:      article.Status,
		LikesCount:  article.LikesCount,
		CreatedAt:   article.CreatedAt,
		UpdatedAt:   article.UpdatedAt,
		PublishedAt: article.PublishedAt,
		User: UserInfo{
			ID:          article.User.ID.String(),
			DisplayName: article.User.DisplayName,
		},
	}

	// シリーズ情報がある場合
	if article.Series != nil {
		detail.Series = &SeriesInfo{
			ID:    article.Series.ID.String(),
			Title: article.Series.Title,
			Slug:  article.Series.Slug,
		}
	}

	// タグ情報の変換
	detail.Tags = make([]TagInfo, len(article.Tags))
	for i, tag := range article.Tags {
		detail.Tags[i] = TagInfo{
			ID:   tag.ID.String(),
			Name: tag.Name,
			Slug: tag.Slug,
		}
	}

	return ArticleDetailResponse{
		Data: detail,
	}
}
