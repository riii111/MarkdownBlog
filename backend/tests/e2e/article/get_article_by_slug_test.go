package article

import (
	"net/http"
	"testing"

	"github.com/riii111/markdown-blog-api/internal/handler/dto"
	"github.com/riii111/markdown-blog-api/tests/e2e"
	"github.com/stretchr/testify/assert"
)

// 記事詳細取得APIのテスト
func TestGetArticleBySlug(t *testing.T) {
	// テスト環境のセットアップ
	router, cleanup := e2e.SetupTestEnvironment(t)
	defer cleanup()

	// テスト用ユーザーの作成とログイン
	_, sessionToken, _ := e2e.CreateAndLoginTestUser(t, router)

	t.Run("正常系: 記事詳細取得成功", func(t *testing.T) {
		// ヘルパー関数を使用して記事作成
		slug := CreateTestArticle(t, router, sessionToken)

		// 記事詳細取得リクエスト
		getURL := "/api/articles/" + slug
		w := e2e.PerformRequest(router, http.MethodGet, getURL, nil, "")

		// ステータスコードの検証
		assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")

		// レスポンスの検証
		var response dto.ArticleDetailResponse
		e2e.GetResponseJSON(t, w, &response)

		assert.NotEmpty(t, response.Data.ID, "Article ID should not be empty")
		assert.Equal(t, slug, response.Data.Slug, "Slug should match")
		assert.NotNil(t, response.Data.User, "User field should not be nil")
	})

	t.Run("異常系: 存在しない記事の詳細取得", func(t *testing.T) {
		// 存在しないslugで記事詳細取得リクエスト
		nonExistentSlug := "non-existent-article-slug"
		getURL := "/api/articles/" + nonExistentSlug
		w := e2e.PerformRequest(router, http.MethodGet, getURL, nil, "")

		// エラーレスポンスの検証
		e2e.AssertErrorResponse(t, w, http.StatusNotFound, "Article not found")
	})

	t.Run("正常系: 複数記事作成後に各記事の詳細取得", func(t *testing.T) {
		// 複数の記事を作成
		slugs := CreateMultipleTestArticles(t, router, sessionToken, 3)

		// 各記事の詳細を取得して検証
		for _, slug := range slugs {
			getURL := "/api/articles/" + slug
			w := e2e.PerformRequest(router, http.MethodGet, getURL, nil, "")

			// ステータスコードの検証
			assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")

			// レスポンスの検証
			var response dto.ArticleDetailResponse
			e2e.GetResponseJSON(t, w, &response)

			assert.NotEmpty(t, response.Data.ID, "Article ID should not be empty")
			assert.Equal(t, slug, response.Data.Slug, "Slug should match")
		}
	})
}
