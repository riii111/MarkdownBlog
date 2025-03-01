package article

import (
	"net/http"
	"testing"

	"github.com/riii111/markdown-blog-api/internal/handler/dto"
	"github.com/riii111/markdown-blog-api/tests/e2e"
	"github.com/stretchr/testify/assert"
)

// 記事削除APIのテスト
func TestDeleteArticle(t *testing.T) {
	// テスト環境のセットアップ
	router, cleanup := e2e.SetupTestEnvironment(t)
	defer cleanup()

	// テスト用ユーザーの作成とログイン
	_, sessionToken, _ := e2e.CreateAndLoginTestUser(t, router)

	t.Run("正常系: 記事削除成功", func(t *testing.T) {
		// 記事作成
		createReq := dto.CreateArticleRequest{}
		w := e2e.PerformRequest(router, http.MethodPost, "/api/articles", createReq, sessionToken)
		assert.Equal(t, http.StatusCreated, w.Code, "Article creation should succeed")

		// 作成した記事のslugを取得
		var createResp dto.CreateArticleResponse
		e2e.GetResponseJSON(t, w, &createResp)
		slug := createResp.Slug

		// 記事削除リクエスト
		deleteURL := "/api/articles/" + slug
		w = e2e.PerformRequest(router, http.MethodDelete, deleteURL, nil, sessionToken)

		// ステータスコードの検証
		assert.Equal(t, http.StatusNoContent, w.Code, "Expected status code 204")
	})

	t.Run("異常系: 認証なしで記事削除", func(t *testing.T) {
		// 記事作成
		createReq := dto.CreateArticleRequest{}
		w := e2e.PerformRequest(router, http.MethodPost, "/api/articles", createReq, sessionToken)
		assert.Equal(t, http.StatusCreated, w.Code, "Article creation should succeed")

		// 作成した記事のslugを取得
		var createResp dto.CreateArticleResponse
		e2e.GetResponseJSON(t, w, &createResp)
		slug := createResp.Slug

		// 認証なしで記事削除リクエスト
		deleteURL := "/api/articles/" + slug
		w = e2e.PerformRequest(router, http.MethodDelete, deleteURL, nil, "")

		// エラーレスポンスの検証
		e2e.AssertErrorResponse(t, w, http.StatusUnauthorized, "Unauthorized")
	})

	t.Run("異常系: 無効なセッショントークンで記事削除", func(t *testing.T) {
		// 記事作成
		createReq := dto.CreateArticleRequest{}
		w := e2e.PerformRequest(router, http.MethodPost, "/api/articles", createReq, sessionToken)
		assert.Equal(t, http.StatusCreated, w.Code, "Article creation should succeed")

		// 作成した記事のslugを取得
		var createResp dto.CreateArticleResponse
		e2e.GetResponseJSON(t, w, &createResp)
		slug := createResp.Slug

		// 無効なセッショントークンで記事削除リクエスト
		deleteURL := "/api/articles/" + slug
		invalidToken := "invalid-session-token"
		w = e2e.PerformRequest(router, http.MethodDelete, deleteURL, nil, invalidToken)

		// エラーレスポンスの検証
		e2e.AssertErrorResponse(t, w, http.StatusUnauthorized, "Unauthorized")
	})

	t.Run("異常系: 存在しない記事を削除", func(t *testing.T) {
		// 存在しないslugで記事削除リクエスト
		nonExistentSlug := "non-existent-article-slug"
		deleteURL := "/api/articles/" + nonExistentSlug
		w := e2e.PerformRequest(router, http.MethodDelete, deleteURL, nil, sessionToken)

		// エラーレスポンスの検証
		e2e.AssertErrorResponse(t, w, http.StatusNotFound, "Article not found")
	})

	t.Run("異常系: 他のユーザーの記事を削除", func(t *testing.T) {
		// 最初のユーザーで記事作成
		createReq := dto.CreateArticleRequest{}
		w := e2e.PerformRequest(router, http.MethodPost, "/api/articles", createReq, sessionToken)
		assert.Equal(t, http.StatusCreated, w.Code, "Article creation should succeed")

		// 作成した記事のslugを取得
		var createResp dto.CreateArticleResponse
		e2e.GetResponseJSON(t, w, &createResp)
		slug := createResp.Slug

		// 別のユーザーを作成してログイン
		_, otherSessionToken, _ := e2e.CreateAndLoginTestUser(t, router)

		// 別のユーザーで最初のユーザーの記事を削除しようとする
		deleteURL := "/api/articles/" + slug
		w = e2e.PerformRequest(router, http.MethodDelete, deleteURL, nil, otherSessionToken)

		// エラーレスポンスの検証（権限がないため削除できない）
		e2e.AssertErrorResponse(t, w, http.StatusUnauthorized, "Unauthorized")
	})
}
