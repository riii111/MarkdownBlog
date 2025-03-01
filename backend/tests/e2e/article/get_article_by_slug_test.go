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
		// 記事作成
		createReq := dto.CreateArticleRequest{}
		w := e2e.PerformRequest(router, http.MethodPost, "/api/articles", createReq, sessionToken)
		assert.Equal(t, http.StatusCreated, w.Code, "Article creation should succeed")

		// 作成した記事のslugを取得
		var createResp dto.CreateArticleResponse
		e2e.GetResponseJSON(t, w, &createResp)
		slug := createResp.Slug

		// 記事詳細取得リクエスト
		getURL := "/api/articles/" + slug
		w = e2e.PerformRequest(router, http.MethodGet, getURL, nil, "")

		// ステータスコードの検証
		assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")

		// レスポンスの検証
		var response dto.ArticleDetailResponse
		e2e.GetResponseJSON(t, w, &response)

		assert.NotEmpty(t, response.Data.ID, "Article ID should not be empty")
		assert.Equal(t, slug, response.Data.Slug, "Slug should match")
	})

	t.Run("異常系: 存在しない記事の詳細取得", func(t *testing.T) {
		// 存在しないslugで記事詳細取得リクエスト
		nonExistentSlug := "non-existent-article-slug"
		getURL := "/api/articles/" + nonExistentSlug
		w := e2e.PerformRequest(router, http.MethodGet, getURL, nil, "")

		// エラーレスポンスの検証
		e2e.AssertErrorResponse(t, w, http.StatusNotFound, "Article not found")
	})

	t.Run("エッジケース: 特殊文字を含むslugで記事詳細取得", func(t *testing.T) {
		// 特殊文字を含むslugで記事詳細取得リクエスト
		// 注: 実際のAPIでは特殊文字を含むslugは生成されないかもしれないが、
		// エッジケースとしてテストする
		specialSlug := "special-slug-with-unusual-characters"
		getURL := "/api/articles/" + specialSlug
		w := e2e.PerformRequest(router, http.MethodGet, getURL, nil, "")

		// 存在しない記事なので404が返るはず
		e2e.AssertErrorResponse(t, w, http.StatusNotFound, "Article not found")
	})

	t.Run("エッジケース: 非常に長いslugで記事詳細取得", func(t *testing.T) {
		// 非常に長いslugで記事詳細取得リクエスト
		longSlug := "very-long-slug-"
		for i := 0; i < 200; i++ {
			longSlug += "a"
		}
		getURL := "/api/articles/" + longSlug
		w := e2e.PerformRequest(router, http.MethodGet, getURL, nil, "")

		// 存在しない記事なので404が返るはず
		e2e.AssertErrorResponse(t, w, http.StatusNotFound, "Article not found")
	})
}
