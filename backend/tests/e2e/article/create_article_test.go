package article

import (
	"net/http"
	"testing"

	"github.com/riii111/markdown-blog-api/internal/handler/dto"
	"github.com/riii111/markdown-blog-api/tests/e2e"
	"github.com/stretchr/testify/assert"
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

	t.Run("正常系: 空のリクエストボディでも記事作成成功", func(t *testing.T) {
		// 空のリクエストボディで記事作成試行
		w := e2e.PerformRequest(router, http.MethodPost, "/api/articles", nil, sessionToken)

		// ステータスコードの検証
		assert.Equal(t, http.StatusCreated, w.Code, "Expected status code 201")

		// レスポンスの検証
		var response dto.CreateArticleResponse
		e2e.GetResponseJSON(t, w, &response)

		assert.NotEmpty(t, response.Slug, "Slug should not be empty")
	})
}
