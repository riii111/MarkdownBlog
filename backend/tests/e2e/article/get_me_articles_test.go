package article

import (
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
	})

	t.Run("正常系: 記事作成後に自分の記事一覧取得", func(t *testing.T) {
		// 記事作成
		createReq := dto.CreateArticleRequest{}
		w := e2e.PerformRequest(router, http.MethodPost, "/api/articles", createReq, sessionToken)
		assert.Equal(t, http.StatusCreated, w.Code, "Article creation should succeed")

		// 自分の記事一覧取得リクエスト
		w = e2e.PerformRequest(router, http.MethodGet, "/api/articles/me", nil, sessionToken)

		// ステータスコードの検証
		assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")

		// レスポンスの検証
		var response dto.ArticleListResponse
		e2e.GetResponseJSON(t, w, &response)

		// 少なくとも1つの記事があるはず
		assert.NotEmpty(t, response.Data, "Articles data should not be empty")
	})

	t.Run("正常系: ページネーションパラメータ指定で自分の記事一覧取得", func(t *testing.T) {
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
		w := e2e.PerformRequest(router, http.MethodGet, "/api/articles/me?page=-1&per_page=0", nil, sessionToken)

		// ステータスコードの検証（APIは不正な値を処理するため200が返るはず）
		assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")

		// レスポンスの検証（デフォルト値が使用されるはず）
		var response dto.ArticleListResponse
		e2e.GetResponseJSON(t, w, &response)

		// デフォルト値（20）が使用されるはず
		assert.Equal(t, 20, response.Pagination.ItemsPerPage, "Items per page should be default (20)")
	})

	t.Run("エッジケース: 過大なページネーションパラメータで自分の記事一覧取得", func(t *testing.T) {
		// 過大なページネーションパラメータで自分の記事一覧取得リクエスト
		w := e2e.PerformRequest(router, http.MethodGet, "/api/articles/me?page=1000&per_page=1000", nil, sessionToken)

		// ステータスコードの検証（APIは過大な値を処理するため200が返るはず）
		assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")

		// レスポンスの検証
		var response dto.ArticleListResponse
		e2e.GetResponseJSON(t, w, &response)

		// 最大値（20）が使用されるはず
		assert.Equal(t, 20, response.Pagination.ItemsPerPage, "Items per page should be max (20)")
	})

	t.Run("複数ユーザーのケース: 他のユーザーの記事は取得されない", func(t *testing.T) {
		// 最初のユーザーで記事作成
		createReq := dto.CreateArticleRequest{}
		w := e2e.PerformRequest(router, http.MethodPost, "/api/articles", createReq, sessionToken)
		assert.Equal(t, http.StatusCreated, w.Code, "Article creation should succeed")

		// 別のユーザーを作成してログイン
		_, otherSessionToken, _ := e2e.CreateAndLoginTestUser(t, router)

		// 別のユーザーで記事作成
		w = e2e.PerformRequest(router, http.MethodPost, "/api/articles", createReq, otherSessionToken)
		assert.Equal(t, http.StatusCreated, w.Code, "Article creation should succeed")

		// 別のユーザーで自分の記事一覧取得リクエスト
		w = e2e.PerformRequest(router, http.MethodGet, "/api/articles/me", nil, otherSessionToken)

		// ステータスコードの検証
		assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")

		// レスポンスの検証
		var response dto.ArticleListResponse
		e2e.GetResponseJSON(t, w, &response)

		// 別のユーザーの記事のみが含まれるはず（1つだけ）
		assert.Equal(t, 1, len(response.Data), "Should only contain the second user's article")
	})
}
