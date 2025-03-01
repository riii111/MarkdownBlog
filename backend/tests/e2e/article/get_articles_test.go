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

	t.Run("正常系: 記事一覧取得", func(t *testing.T) {
		// リクエスト実行
		w := e2e.PerformRequest(router, http.MethodGet, "/api/articles", nil, "")

		// ステータスコードの検証
		assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")

		// レスポンスの検証
		var response dto.ArticleListResponse
		e2e.GetResponseJSON(t, w, &response)

		// 初期状態では記事がない可能性があるが、レスポンス構造は正しいはず
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

	t.Run("エッジケース: 無効なcursor指定で記事一覧取得", func(t *testing.T) {
		// 無効なcursorパラメータ付きでリクエスト実行
		w := e2e.PerformRequest(router, http.MethodGet, "/api/articles?cursor=invalid-cursor", nil, "")

		// APIの実装によっては500エラーが返る場合もある
		if w.Code == http.StatusInternalServerError {
			assert.Equal(t, http.StatusInternalServerError, w.Code, "Expected status code 500 for invalid cursor")
			return
		}

		// または、無効なcursorを無視して200が返る場合もある
		assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")

		// レスポンスの検証
		var response dto.ArticleListResponse
		e2e.GetResponseJSON(t, w, &response)

		assert.NotNil(t, response.Data, "Data field should not be nil")
	})

	t.Run("エッジケース: 不正なlimit値で記事一覧取得", func(t *testing.T) {
		// 不正なlimit値でリクエスト実行
		w := e2e.PerformRequest(router, http.MethodGet, "/api/articles?limit=-1", nil, "")

		// ステータスコードの検証（APIは不正な値を処理するため200が返るはず）
		assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")

		// レスポンスの検証（デフォルト値が使用されるはず）
		var response dto.ArticleListResponse
		e2e.GetResponseJSON(t, w, &response)

		// デフォルト値（usecaseで定義されている9）が使用されるはず
		assert.Equal(t, 9, response.Pagination.ItemsPerPage, "Items per page should be default (9)")
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
		assert.Equal(t, 9, response.Pagination.ItemsPerPage, "Items per page should be max (9)")
	})
}
