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
