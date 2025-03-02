package article

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/riii111/markdown-blog-api/internal/handler/dto"
	"github.com/riii111/markdown-blog-api/tests/e2e"
	"github.com/stretchr/testify/assert"
)

// 自分の記事一覧取得APIのテスト
func TestGetMeArticles(t *testing.T) {
	// テスト環境のセットアップ
	router, cleanup := e2e.SetupTestEnvironment(t)
	defer cleanup()

	// テスト用ユーザーの作成とログイン
	_, sessionToken, _ := e2e.CreateAndLoginTestUser(t, router)

	t.Run("正常系: 自分の記事一覧取得（記事なし）", func(t *testing.T) {
		// 自分の記事一覧取得リクエスト
		w := e2e.PerformRequest(router, http.MethodGet, "/api/articles/me", nil, sessionToken)

		// ステータスコードの検証
		assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")

		// レスポンスの検証
		var response dto.ArticleListResponse
		e2e.GetResponseJSON(t, w, &response)

		// 初期状態では記事がないはず
		assert.NotNil(t, response.Data, "Data field should not be nil")
		assert.NotNil(t, response.Pagination, "Pagination field should not be nil")
		assert.Empty(t, response.Data, "Articles data should be empty")
	})

	t.Run("正常系: 記事作成後に自分の記事一覧取得", func(t *testing.T) {
		// 記事作成
		slug := CreateTestArticle(t, router, sessionToken)
		assert.NotEmpty(t, slug, "Article creation should succeed")

		// 自分の記事一覧取得リクエスト
		w := e2e.PerformRequest(router, http.MethodGet, "/api/articles/me", nil, sessionToken)

		// ステータスコードの検証
		assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")

		// レスポンスの検証
		var response dto.ArticleListResponse
		e2e.GetResponseJSON(t, w, &response)

		// 少なくとも1つの記事があるはず
		assert.NotEmpty(t, response.Data, "Articles data should not be empty")
	})

	t.Run("正常系: ページネーションパラメータ指定で自分の記事一覧取得", func(t *testing.T) {
		// 複数の記事を作成
		CreateMultipleTestArticles(t, router, sessionToken, 3)

		// ページネーションパラメータ付きで自分の記事一覧取得リクエスト
		w := e2e.PerformRequest(router, http.MethodGet, "/api/articles/me?page=1&per_page=5", nil, sessionToken)

		// ステータスコードの検証
		assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")

		// レスポンスの検証
		var response dto.ArticleListResponse
		e2e.GetResponseJSON(t, w, &response)

		assert.NotNil(t, response.Data, "Data field should not be nil")
		assert.Equal(t, 5, response.Pagination.ItemsPerPage, "Items per page should be 5")
	})

	t.Run("異常系: 認証なしで自分の記事一覧取得", func(t *testing.T) {
		// 認証なしで自分の記事一覧取得リクエスト
		w := e2e.PerformRequest(router, http.MethodGet, "/api/articles/me", nil, "")

		// エラーレスポンスの検証
		e2e.AssertErrorResponse(t, w, http.StatusUnauthorized, "Unauthorized")
	})

	t.Run("異常系: 無効なセッショントークンで自分の記事一覧取得", func(t *testing.T) {
		// 無効なセッショントークンで自分の記事一覧取得リクエスト
		invalidToken := "invalid-session-token"
		w := e2e.PerformRequest(router, http.MethodGet, "/api/articles/me", nil, invalidToken)

		// エラーレスポンスの検証
		e2e.AssertErrorResponse(t, w, http.StatusUnauthorized, "Unauthorized")
	})

	t.Run("エッジケース: 不正なページネーションパラメータで自分の記事一覧取得", func(t *testing.T) {
		// 不正なページネーションパラメータで自分の記事一覧取得リクエスト
		w := e2e.PerformRequest(router, http.MethodGet,
			fmt.Sprintf("/api/articles/me?page=%d&per_page=%d", InvalidValue, InvalidValue),
			nil, sessionToken)

		// ステータスコードの検証（APIは不正な値を処理するため200が返るはず）
		assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")

		// レスポンスの検証（デフォルト値が使用されるはず）
		var response dto.ArticleListResponse
		e2e.GetResponseJSON(t, w, &response)

		// デフォルト値（DefaultPageSize）が使用されるはず
		assert.Equal(t, DefaultPageSize, response.Pagination.ItemsPerPage,
			fmt.Sprintf("Items per page should be default (%d)", DefaultPageSize))
	})

	t.Run("エッジケース: 過大なページネーションパラメータで自分の記事一覧取得", func(t *testing.T) {
		// 過大なページネーションパラメータで自分の記事一覧取得リクエスト
		w := e2e.PerformRequest(router, http.MethodGet,
			fmt.Sprintf("/api/articles/me?page=%d&per_page=%d", ExcessiveValue, ExcessiveValue),
			nil, sessionToken)

		// ステータスコードの検証（APIは過大な値を処理するため200が返るはず）
		assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")

		// レスポンスの検証
		var response dto.ArticleListResponse
		e2e.GetResponseJSON(t, w, &response)

		// 最大値（DefaultPageSize）が使用されるはず
		assert.Equal(t, DefaultPageSize, response.Pagination.ItemsPerPage,
			fmt.Sprintf("Items per page should be max (%d)", DefaultPageSize))
	})

	t.Run("複数ユーザーのケース: 他のユーザーの記事は取得されない", func(t *testing.T) {
		// 新しいテスト環境をセットアップして、前のテストの影響を受けないようにする
		newRouter, newCleanup := e2e.SetupTestEnvironment(t)
		defer newCleanup()

		// 最初のユーザーを作成してログイン
		_, firstUserToken, _ := e2e.CreateAndLoginTestUser(t, newRouter)

		// 最初のユーザーで記事作成
		firstUserSlug := CreateTestArticle(t, newRouter, firstUserToken)
		assert.NotEmpty(t, firstUserSlug, "First user's article creation should succeed")

		// 別のユーザーを作成してログイン
		_, secondUserToken, _ := e2e.CreateAndLoginTestUser(t, newRouter)

		// 別のユーザーで記事作成
		secondUserSlug := CreateTestArticle(t, newRouter, secondUserToken)
		assert.NotEmpty(t, secondUserSlug, "Second user's article creation should succeed")

		// 別のユーザーで自分の記事一覧取得リクエスト
		w := e2e.PerformRequest(newRouter, http.MethodGet, "/api/articles/me", nil, secondUserToken)

		// ステータスコードの検証
		assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")

		// レスポンスの検証
		var response dto.ArticleListResponse
		e2e.GetResponseJSON(t, w, &response)

		// 別のユーザーの記事のみが含まれるはず（1つだけ）
		assert.Equal(t, 1, len(response.Data), "Should only contain the second user's article")

		// 最初のユーザーで自分の記事一覧取得リクエスト
		w = e2e.PerformRequest(newRouter, http.MethodGet, "/api/articles/me", nil, firstUserToken)

		// レスポンスの検証
		var firstUserResponse dto.ArticleListResponse
		e2e.GetResponseJSON(t, w, &firstUserResponse)

		// 最初のユーザーの記事のみが含まれるはず（1つだけ）
		assert.Equal(t, 1, len(firstUserResponse.Data), "Should only contain the first user's article")
	})
}
