package article

import (
	"net/http"
	"testing"

	"github.com/riii111/markdown-blog-api/internal/handler/dto"
	"github.com/riii111/markdown-blog-api/tests/e2e"
	"github.com/stretchr/testify/assert"
)

// 公開記事一覧取得APIのテスト
func TestGetArticles(t *testing.T) {
	// テスト環境のセットアップ
	router, cleanup := e2e.SetupTestEnvironment(t)
	defer cleanup()

	// テスト用ユーザーの作成とログイン
	_, sessionToken, _ := e2e.CreateAndLoginTestUser(t, router)

	CreateTestArticle(t, router, sessionToken)

	t.Run("正常系: 記事一覧取得", func(t *testing.T) {
		// リクエスト実行
		w := e2e.PerformRequest(router, http.MethodGet, "/api/articles", nil, "")

		// ステータスコードの検証
		assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")

		// レスポンスの検証
		var response dto.ArticleListResponse
		e2e.GetResponseJSON(t, w, &response)

		assert.NotNil(t, response.Data, "Data field should not be nil")
		assert.NotNil(t, response.Pagination, "Pagination field should not be nil")
	})

	t.Run("正常系: limit指定で記事一覧取得", func(t *testing.T) {
		// limitパラメータ付きでリクエスト実行
		w := e2e.PerformRequest(router, http.MethodGet, "/api/articles?limit=5", nil, "")

		// ステータスコードの検証
		assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")

		// レスポンスの検証
		var response dto.ArticleListResponse
		e2e.GetResponseJSON(t, w, &response)

		assert.Equal(t, 5, response.Pagination.ItemsPerPage, "Items per page should be 5")
	})

	t.Run("正常系: 複数記事作成後の一覧取得", func(t *testing.T) {
		// 複数の記事を作成
		CreateMultipleTestArticles(t, router, sessionToken, 3)

		// リクエスト実行
		w := e2e.PerformRequest(router, http.MethodGet, "/api/articles", nil, "")

		// ステータスコードの検証
		assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")

		// レスポンスの検証
		var response dto.ArticleListResponse
		e2e.GetResponseJSON(t, w, &response)

		// 注: 作成した記事はドラフト状態のため公開記事一覧には含まれない
		assert.NotNil(t, response.Data, "Data field should not be nil")
	})

	t.Run("エッジケース: 無効なcursor指定で記事一覧取得", func(t *testing.T) {
		// 無効なcursorパラメータ付きでリクエスト実行
		w := e2e.PerformRequest(router, http.MethodGet, "/api/articles?cursor=invalid-cursor", nil, "")

		// 無効なcursorの場合は500エラーが返される
		if w.Code == http.StatusInternalServerError {
			// 500エラーの場合はエラーレスポンスを検証
			e2e.AssertErrorResponse(t, w, http.StatusInternalServerError, "Failed to fetch articles")
			return
		}

		// このケースは発生しないはずだが、念のため検証
		t.Errorf("無効なcursorの場合は500エラーが返されるはずですが、200が返されました")
	})

	// 境界値テスト
	t.Run("エッジケース: 不正なlimit値で記事一覧取得", func(t *testing.T) {
		// 不正なlimit値でリクエスト実行
		w := e2e.PerformRequest(router, http.MethodGet, "/api/articles?limit=-1", nil, "")

		// ステータスコードの検証（APIは不正な値を処理するため200が返るはず）
		assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")

		// レスポンスの検証（デフォルト値が使用されるはず）
		var response dto.ArticleListResponse
		e2e.GetResponseJSON(t, w, &response)

		// デフォルト値（usecaseで定義されている9）が使用されるはず
		assert.Equal(t, DefaultArticleLimit, response.Pagination.ItemsPerPage, "Items per page should be default (9)")
	})

	t.Run("エッジケース: 過大なlimit値で記事一覧取得", func(t *testing.T) {
		// 過大なlimit値でリクエスト実行
		w := e2e.PerformRequest(router, http.MethodGet, "/api/articles?limit=1000", nil, "")

		// ステータスコードの検証（APIは過大な値を処理するため200が返るはず）
		assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")

		// レスポンスの検証（最大値が使用されるはず）
		var response dto.ArticleListResponse
		e2e.GetResponseJSON(t, w, &response)

		// 最大値（usecaseで定義されている9）が使用されるはず
		assert.Equal(t, DefaultArticleLimit, response.Pagination.ItemsPerPage, "Items per page should be max (9)")
	})

	t.Run("エッジケース: 複数パラメータ指定で記事一覧取得", func(t *testing.T) {
		// 複数パラメータ付きでリクエスト実行
		w := e2e.PerformRequest(router, http.MethodGet, "/api/articles?limit=3&cursor=test", nil, "")

		// 無効なcursorの場合は500エラーが返される
		if w.Code == http.StatusInternalServerError {
			// 500エラーの場合はエラーレスポンスを検証
			e2e.AssertErrorResponse(t, w, http.StatusInternalServerError, "Failed to fetch articles")
			return
		}

		// 正常に処理された場合
		assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")

		// レスポンスの検証
		var response dto.ArticleListResponse
		e2e.GetResponseJSON(t, w, &response)

		assert.NotNil(t, response.Data, "Data field should not be nil")
		assert.Equal(t, 3, response.Pagination.ItemsPerPage, "Items per page should be 3")
	})
}
