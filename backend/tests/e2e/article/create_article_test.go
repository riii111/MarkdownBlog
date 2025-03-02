package article

import (
	"net/http"
	"testing"

	"github.com/riii111/markdown-blog-api/internal/handler/dto"
	"github.com/riii111/markdown-blog-api/tests/e2e"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 記事作成APIのテスト
func TestCreateArticle(t *testing.T) {
	// テスト環境のセットアップ
	router, cleanup := e2e.SetupTestEnvironment(t)
	defer cleanup()

	// テスト用ユーザーの作成とログイン
	_, sessionToken, _ := e2e.CreateAndLoginTestUser(t, router)

	t.Run("正常系: 記事作成成功", func(t *testing.T) {
		// 記事作成リクエスト（空の構造体）
		createReq := dto.CreateArticleRequest{}

		// リクエスト実行
		w := e2e.PerformRequest(router, http.MethodPost, "/api/articles", createReq, sessionToken)

		// ステータスコードの検証
		assert.Equal(t, http.StatusCreated, w.Code, "Expected status code 201")

		// レスポンスの検証
		var response dto.CreateArticleResponse
		e2e.GetResponseJSON(t, w, &response)

		assert.NotEmpty(t, response.Slug, "Slug should not be empty")

		// 作成された記事の取得と検証
		w = e2e.PerformRequest(router, http.MethodGet, "/api/articles/"+response.Slug, nil, sessionToken)
		assert.Equal(t, http.StatusOK, w.Code, "Should be able to get the created article")

		var articleResponse dto.ArticleDetailResponse
		e2e.GetResponseJSON(t, w, &articleResponse)
		require.NotNil(t, articleResponse.Data, "Article data should not be nil")
		assert.Equal(t, response.Slug, articleResponse.Data.Slug, "Slug should match")
		assert.NotNil(t, articleResponse.Data.User, "User field should not be nil")
	})

	t.Run("正常系: 空のリクエストボディでも記事作成成功", func(t *testing.T) {
		// 空のリクエストボディで記事作成試行
		w := e2e.PerformRequest(router, http.MethodPost, "/api/articles", nil, sessionToken)

		// ステータスコードの検証
		assert.Equal(t, http.StatusCreated, w.Code, "Expected status code 201")

		// レスポンスの検証
		var response dto.CreateArticleResponse
		e2e.GetResponseJSON(t, w, &response)

		assert.NotEmpty(t, response.Slug, "Slug should not be empty")

		// 作成された記事の取得と検証
		w = e2e.PerformRequest(router, http.MethodGet, "/api/articles/"+response.Slug, nil, sessionToken)
		assert.Equal(t, http.StatusOK, w.Code, "Should be able to get the created article")
	})

	t.Run("正常系: 連続して複数の記事を作成できる", func(t *testing.T) {
		// 複数の記事を作成
		slugs := CreateMultipleTestArticles(t, router, sessionToken, 3)

		// 各記事が取得できることを確認
		for _, slug := range slugs {
			w := e2e.PerformRequest(router, http.MethodGet, "/api/articles/"+slug, nil, sessionToken)
			assert.Equal(t, http.StatusOK, w.Code, "Should be able to get the created article")
		}
	})

	t.Run("異常系: 認証なしで記事作成", func(t *testing.T) {
		// 認証なしで記事作成リクエスト
		createReq := dto.CreateArticleRequest{}

		// リクエスト実行（セッショントークンなし）
		w := e2e.PerformRequest(router, http.MethodPost, "/api/articles", createReq, "")

		// エラーレスポンスの検証
		e2e.AssertErrorResponse(t, w, http.StatusUnauthorized, "Unauthorized")
	})

	t.Run("異常系: 無効なセッショントークンで記事作成", func(t *testing.T) {
		// 無効なセッショントークンで記事作成リクエスト
		createReq := dto.CreateArticleRequest{}

		// リクエスト実行（無効なセッショントークン）
		invalidToken := "invalid-session-token"
		w := e2e.PerformRequest(router, http.MethodPost, "/api/articles", createReq, invalidToken)

		// エラーレスポンスの検証
		e2e.AssertErrorResponse(t, w, http.StatusUnauthorized, "Unauthorized")
	})
}
