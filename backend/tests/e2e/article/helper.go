package article

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/riii111/markdown-blog-api/internal/handler/dto"
	"github.com/riii111/markdown-blog-api/tests/e2e"
	"github.com/stretchr/testify/require"
)

// ページネーションとパラメータの定数
// APIの仕様変更に対応しやすくするため、定数として定義
const (
	InvalidValue        = -1
	ExcessiveValue      = 1000
	DefaultPageSize     = 20
	MaxPageSize         = 100
	DefaultArticleLimit = 9
	MaxArticleLimit     = 100
)

// 記事作成のヘルパー関数
// テスト用の記事を作成し、そのslugを返す
func CreateTestArticle(t *testing.T, router *gin.Engine, sessionToken string) string {
	createReq := dto.CreateArticleRequest{}
	w := e2e.PerformRequest(router, http.MethodPost, "/api/articles", createReq, sessionToken)
	require.Equal(t, http.StatusCreated, w.Code, "Article creation should succeed")

	var createResp dto.CreateArticleResponse
	e2e.GetResponseJSON(t, w, &createResp)
	return createResp.Slug
}

// 複数の記事を作成するヘルパー関数
// 指定された数の記事を作成し、そのslugのスライスを返す
func CreateMultipleTestArticles(t *testing.T, router *gin.Engine, sessionToken string, count int) []string {
	slugs := make([]string, count)
	for i := 0; i < count; i++ {
		slugs[i] = CreateTestArticle(t, router, sessionToken)
	}
	return slugs
}

// 記事のステータスを公開状態に変更するヘルパー関数
// 記事のslugを受け取り、その記事を公開状態に変更する
func PublishArticle(t *testing.T, router *gin.Engine, sessionToken string, slug string) {
	// 記事更新用のリクエスト
	updateReq := dto.UpdateArticleRequest{
		Status: "published",
	}
	
	// 記事更新リクエストを実行
	updateURL := fmt.Sprintf("/api/articles/%s", slug)
	w := e2e.PerformRequest(router, http.MethodPut, updateURL, updateReq, sessionToken)
	require.Equal(t, http.StatusOK, w.Code, "Article status update should succeed")
}

// 複数の記事を公開状態に変更するヘルパー関数
func PublishMultipleArticles(t *testing.T, router *gin.Engine, sessionToken string, slugs []string) {
	for _, slug := range slugs {
		PublishArticle(t, router, sessionToken, slug)
	}
}

// テスト用のタグを作成するヘルパー関数
// このテストではダミーのタグスラグを使用します
func CreateTestTag(t *testing.T, router *gin.Engine, name string) string {
	// ダミーのスラグを返します
	return fmt.Sprintf("test-tag-%s", uuid.New().String()[:8])
}

// テスト用のタグ付き記事を作成するヘルパー関数
// このテストではタグ付き記事のシミュレーションを行います
func CreateArticleWithTag(t *testing.T, router *gin.Engine, sessionToken string, tagSlug string) string {
	// 記事作成
	slug := CreateTestArticle(t, router, sessionToken)
	return slug
}
