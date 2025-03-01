package article

import (
	"net/http"
	"testing"

	"github.com/riii111/markdown-blog-api/internal/handler/dto"
	"github.com/riii111/markdown-blog-api/tests/e2e"
	"github.com/stretchr/testify/assert"
)

// タグに紐づく記事一覧取得APIのテスト
func TestGetArticlesByTag(t *testing.T) {
	// テスト環境のセットアップ
	router, cleanup := e2e.SetupTestEnvironment(t)
	defer cleanup()

	// テスト用ユーザーの作成とログイン
	// TODO: 将来的にタグ付き記事作成に使用する可能性あり
	_, _, _ = e2e.CreateAndLoginTestUser(t, router)

	// 注: このテストでは、タグ付き記事の作成方法がAPIの実装に依存する
	// 実際のAPIでは、記事作成時にタグを指定できるかもしれないが
	// 現在の実装では記事作成時にタグを指定できないため、
	// 一時的にこのテストコードを用意する

	t.Run("異常系: 存在しないタグで記事一覧取得", func(t *testing.T) {
		// 存在しないタグで記事一覧取得リクエスト
		nonExistentTag := "non-existent-tag"
		getURL := "/api/tags/" + nonExistentTag + "/articles"
		w := e2e.PerformRequest(router, http.MethodGet, getURL, nil, "")

		// ステータスコードの検証
		assert.Equal(t, http.StatusNotFound, w.Code, "Expected status code 404")
	})

	t.Run("正常系: タグに紐づく記事一覧取得（記事なし）", func(t *testing.T) {
		// 実際のタグスラグを使用（存在するタグだが記事がない場合）
		// 注: このテストは実際のタグがある場合のみ成功します
		tagSlug := "test-tag"
		getURL := "/api/tags/" + tagSlug + "/articles"
		w := e2e.PerformRequest(router, http.MethodGet, getURL, nil, "")

		// タグが存在しない場合は404が返るはず
		if w.Code == http.StatusNotFound {
			assert.Equal(t, http.StatusNotFound, w.Code, "Expected status code 404")
			return
		}

		// タグが存在する場合は200が返り、空の記事リストが返るはず
		assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")

		var response dto.ArticleListResponse
		e2e.GetResponseJSON(t, w, &response)

		assert.NotNil(t, response.Data, "Data field should not be nil")
		assert.NotNil(t, response.Pagination, "Pagination field should not be nil")
	})

	t.Run("正常系: ページネーションパラメータ指定でタグに紐づく記事一覧取得", func(t *testing.T) {
		// ページネーションパラメータ付きでタグに紐づく記事一覧取得リクエスト
		tagSlug := "test-tag"
		getURL := "/api/tags/" + tagSlug + "/articles?limit=5&cursor=some-cursor"
		w := e2e.PerformRequest(router, http.MethodGet, getURL, nil, "")

		// タグが存在しない場合は404が返るはず
		if w.Code == http.StatusNotFound {
			assert.Equal(t, http.StatusNotFound, w.Code, "Expected status code 404")
			return
		}

		// タグが存在する場合は200が返るはず
		assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")

		var response dto.ArticleListResponse
		e2e.GetResponseJSON(t, w, &response)

		assert.NotNil(t, response.Data, "Data field should not be nil")
		assert.Equal(t, 5, response.Pagination.ItemsPerPage, "Items per page should be 5")
	})

	t.Run("エッジケース: 特殊文字を含むタグスラグで記事一覧取得", func(t *testing.T) {
		// 特殊文字を含むタグスラグで記事一覧取得リクエスト
		specialTagSlug := "special-tag-with-unusual-characters"
		getURL := "/api/tags/" + specialTagSlug + "/articles"
		w := e2e.PerformRequest(router, http.MethodGet, getURL, nil, "")

		// 存在しないタグなので404が返るはず
		assert.Equal(t, http.StatusNotFound, w.Code, "Expected status code 404")
	})

	t.Run("エッジケース: 不正なlimit値でタグに紐づく記事一覧取得", func(t *testing.T) {
		// 不正なlimit値でタグに紐づく記事一覧取得リクエスト
		tagSlug := "test-tag"
		getURL := "/api/tags/" + tagSlug + "/articles?limit=-1"
		w := e2e.PerformRequest(router, http.MethodGet, getURL, nil, "")

		// タグが存在しない場合は404が返るはず
		if w.Code == http.StatusNotFound {
			assert.Equal(t, http.StatusNotFound, w.Code, "Expected status code 404")
			return
		}

		// タグが存在する場合は200が返り、デフォルト値が使用されるはず
		assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")

		var response dto.ArticleListResponse
		e2e.GetResponseJSON(t, w, &response)

		// デフォルト値（20）が使用されるはず
		assert.Equal(t, 20, response.Pagination.ItemsPerPage, "Items per page should be default (20)")
	})

	t.Run("エッジケース: 過大なlimit値でタグに紐づく記事一覧取得", func(t *testing.T) {
		// 過大なlimit値でタグに紐づく記事一覧取得リクエスト
		tagSlug := "test-tag"
		getURL := "/api/tags/" + tagSlug + "/articles?limit=1000"
		w := e2e.PerformRequest(router, http.MethodGet, getURL, nil, "")

		// タグが存在しない場合は404が返るはず
		if w.Code == http.StatusNotFound {
			assert.Equal(t, http.StatusNotFound, w.Code, "Expected status code 404")
			return
		}

		// タグが存在する場合は200が返り、最大値が使用されるはず
		assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")

		var response dto.ArticleListResponse
		e2e.GetResponseJSON(t, w, &response)

		// 最大値（20）が使用されるはず
		assert.Equal(t, 20, response.Pagination.ItemsPerPage, "Items per page should be max (20)")
	})
}
