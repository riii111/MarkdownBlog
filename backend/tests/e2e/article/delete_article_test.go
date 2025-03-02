package article

import (
	"net/http"
	"testing"

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

	// テスト環境はSetupTestEnvironmentのcleanup関数によってクリーンアップされる

	t.Run("正常系: 記事削除成功", func(t *testing.T) {
		// ヘルパー関数を使用して記事作成
		slug := CreateTestArticle(t, router, sessionToken)

		// 記事削除リクエスト
		deleteURL := "/api/articles/" + slug
		w := e2e.PerformRequest(router, http.MethodDelete, deleteURL, nil, sessionToken)

		// ステータスコードの検証
		assert.Equal(t, http.StatusNoContent, w.Code, "Expected status code 204")

		// 削除後に記事が存在しないことを確認
		getURL := "/api/articles/" + slug
		w = e2e.PerformRequest(router, http.MethodGet, getURL, nil, "")
		e2e.AssertErrorResponse(t, w, http.StatusNotFound, "Article not found")
	})

	t.Run("異常系: 認証なしで記事削除", func(t *testing.T) {
		// ヘルパー関数を使用して記事作成
		slug := CreateTestArticle(t, router, sessionToken)

		// 認証なしで記事削除リクエスト
		deleteURL := "/api/articles/" + slug
		w := e2e.PerformRequest(router, http.MethodDelete, deleteURL, nil, "")

		// エラーレスポンスの検証
		e2e.AssertErrorResponse(t, w, http.StatusUnauthorized, "Unauthorized")

		// 記事が削除されていないことを確認
		getURL := "/api/articles/" + slug
		w = e2e.PerformRequest(router, http.MethodGet, getURL, nil, "")
		assert.Equal(t, http.StatusOK, w.Code, "Article should still exist")
	})

	t.Run("異常系: 無効なセッショントークンで記事削除", func(t *testing.T) {
		// ヘルパー関数を使用して記事作成
		slug := CreateTestArticle(t, router, sessionToken)

		// 無効なセッショントークンで記事削除リクエスト
		deleteURL := "/api/articles/" + slug
		invalidToken := "invalid-session-token"
		w := e2e.PerformRequest(router, http.MethodDelete, deleteURL, nil, invalidToken)

		// エラーレスポンスの検証
		e2e.AssertErrorResponse(t, w, http.StatusUnauthorized, "Unauthorized")

		// 記事が削除されていないことを確認
		getURL := "/api/articles/" + slug
		w = e2e.PerformRequest(router, http.MethodGet, getURL, nil, "")
		assert.Equal(t, http.StatusOK, w.Code, "Article should still exist")
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
		slug := CreateTestArticle(t, router, sessionToken)

		// 別のユーザーを作成してログイン
		_, otherSessionToken, _ := e2e.CreateAndLoginTestUser(t, router)

		// 別のユーザーで最初のユーザーの記事を削除しようとする
		deleteURL := "/api/articles/" + slug
		w := e2e.PerformRequest(router, http.MethodDelete, deleteURL, nil, otherSessionToken)

		// エラーレスポンスの検証（権限がないため削除できない）
		e2e.AssertErrorResponse(t, w, http.StatusUnauthorized, "Unauthorized")

		// 記事が削除されていないことを確認
		getURL := "/api/articles/" + slug
		w = e2e.PerformRequest(router, http.MethodGet, getURL, nil, "")
		assert.Equal(t, http.StatusOK, w.Code, "Article should still exist")
	})

	t.Run("エッジケース: 削除済みの記事を再度削除", func(t *testing.T) {
		// ヘルパー関数を使用して記事作成
		slug := CreateTestArticle(t, router, sessionToken)

		// 記事削除リクエスト
		deleteURL := "/api/articles/" + slug
		w := e2e.PerformRequest(router, http.MethodDelete, deleteURL, nil, sessionToken)
		assert.Equal(t, http.StatusNoContent, w.Code, "First deletion should succeed")

		// 同じ記事を再度削除
		w = e2e.PerformRequest(router, http.MethodDelete, deleteURL, nil, sessionToken)

		// 記事が既に存在しないため404が返るはず
		e2e.AssertErrorResponse(t, w, http.StatusNotFound, "Article not found")
	})
}
