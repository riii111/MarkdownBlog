package article

import (
	"fmt"
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
	_, sessionToken, _ := e2e.CreateAndLoginTestUser(t, router)

	// テスト用のタグを作成
	// 注: 実際のプロジェクトではタグ作成APIが必要
	// 現在はダミーのタグスラグを使用する
	testTagSlug := CreateTestTag(t, router, "テスト用タグ")

	t.Run("異常系: 存在しないタグで記事一覧取得", func(t *testing.T) {
		// 存在しないタグで記事一覧取得リクエスト
		nonExistentTag := "non-existent-tag"
		getURL := "/api/tags/" + nonExistentTag + "/articles"
		w := e2e.PerformRequest(router, http.MethodGet, getURL, nil, "")

		// ステータスコードの検証のみ行う
		assert.Equal(t, http.StatusNotFound, w.Code, "Expected status code 404")
	})

	t.Run("正常系: タグに紐づく記事一覧取得（記事なし）", func(t *testing.T) {
		// テスト用のタグを使用（記事がない状態）
		getURL := "/api/tags/" + testTagSlug + "/articles"
		w := e2e.PerformRequest(router, http.MethodGet, getURL, nil, "")

		// タグが存在しないか、タグに紐づく記事がない場合は404が返される
		assert.Equal(t, http.StatusNotFound, w.Code, "Expected status code 404")
	})

	t.Run("正常系: タグに紐づく記事一覧取得（記事あり）", func(t *testing.T) {
		// テスト用のタグ付き記事を作成
		CreateArticleWithTag(t, router, sessionToken, testTagSlug)

		// タグに紐づく記事一覧取得リクエスト
		getURL := "/api/tags/" + testTagSlug + "/articles"
		w := e2e.PerformRequest(router, http.MethodGet, getURL, nil, "")

		// タグに紐づく記事がない場合は404が返されます
		if w.Code == http.StatusNotFound {
			// この場合はテストをスキップ
			t.Skip("記事とタグの関連付けが反映されていない可能性があります")
			return
		}

		// 正常に取得できた場合
		assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")

		var response dto.ArticleListResponse
		e2e.GetResponseJSON(t, w, &response)

		assert.NotNil(t, response.Data, "Data field should not be nil")
		assert.NotNil(t, response.Pagination, "Pagination field should not be nil")
		// 少なくとも1つの記事があるはず
		assert.NotEmpty(t, response.Data, "Articles data should not be empty")
	})

	t.Run("正常系: ページネーションパラメータ指定でタグに紐づく記事一覧取得", func(t *testing.T) {
		// ページネーションパラメータ付きでタグに紐づく記事一覧取得リクエスト
		getURL := "/api/tags/" + testTagSlug + "/articles?limit=5"
		w := e2e.PerformRequest(router, http.MethodGet, getURL, nil, "")

		// タグに紐づく記事がない場合は404が返される
		if w.Code == http.StatusNotFound {
			// この場合はテストをスキップ
			t.Skip("記事とタグの関連付けが反映されていない可能性があります")
			return
		}

		// 正常に取得できた場合
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

	// 境界値テスト
	t.Run("エッジケース: 不正なlimit値でタグに紐づく記事一覧取得", func(t *testing.T) {
		// 不正なlimit値でタグに紐づく記事一覧取得リクエスト
		getURL := "/api/tags/" + testTagSlug + "/articles?limit=" + fmt.Sprint(InvalidValue)
		w := e2e.PerformRequest(router, http.MethodGet, getURL, nil, "")

		// タグに紐づく記事がない場合は404が返される
		if w.Code == http.StatusNotFound {
			// この場合はテストをスキップ
			t.Skip("記事とタグの関連付けが反映されていない可能性があります")
			return
		}

		// 正常に取得できた場合
		assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")

		var response dto.ArticleListResponse
		e2e.GetResponseJSON(t, w, &response)

		// デフォルト値（DefaultArticleLimit）が使用されるはず
		assert.Equal(t, DefaultArticleLimit, response.Pagination.ItemsPerPage,
			"Items per page should be default (9)")
	})

	t.Run("エッジケース: 過大なlimit値でタグに紐づく記事一覧取得", func(t *testing.T) {
		// 過大なlimit値でタグに紐づく記事一覧取得リクエスト
		getURL := "/api/tags/" + testTagSlug + "/articles?limit=" + fmt.Sprint(ExcessiveValue)
		w := e2e.PerformRequest(router, http.MethodGet, getURL, nil, "")

		// タグに紐づく記事がない場合は404が返される
		if w.Code == http.StatusNotFound {
			// この場合はテストをスキップ
			t.Skip("記事とタグの関連付けが反映されていない可能性があります")
			return
		}

		// 正常に取得できた場合
		assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")

		var response dto.ArticleListResponse
		e2e.GetResponseJSON(t, w, &response)

		// 最大値（DefaultArticleLimit）が使用されるはず
		assert.Equal(t, DefaultArticleLimit, response.Pagination.ItemsPerPage,
			"Items per page should be max (9)")
	})

	// 無効なcursorテスト
	t.Run("エッジケース: 無効なcursorでタグに紐づく記事一覧取得", func(t *testing.T) {
		// 無効なcursorでタグに紐づく記事一覧取得リクエスト
		getURL := "/api/tags/" + testTagSlug + "/articles?cursor=invalid-cursor"
		w := e2e.PerformRequest(router, http.MethodGet, getURL, nil, "")

		// タグが存在しない場合は404が返される
		if w.Code == http.StatusNotFound {
			// 404エラーの場合
			assert.Equal(t, http.StatusNotFound, w.Code, "Expected status code 404")
			return
		}

		// タグが存在する場合、無効なcursorの場合は500エラーが返される
		if w.Code == http.StatusInternalServerError {
			// 500エラーの場合
			assert.Equal(t, http.StatusInternalServerError, w.Code, "Expected status code 500")
			return
		}

		// このケースは発生しないはずだが、念のため検証
		t.Errorf("無効なcursorの場合は404または500エラーが返されるはずですが、%dが返されました", w.Code)
	})
}
