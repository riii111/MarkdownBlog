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

	// 記事を作成して公開状態にする
	slug := CreateTestArticle(t, router, sessionToken)
	PublishArticle(t, router, sessionToken, slug)

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
		// 公開記事が少なくとも1つあるはず
		assert.NotEmpty(t, response.Data, "Articles data should not be empty")
		
		// すべての記事が公開状態であることを検証
		for _, article := range response.Data {
			assert.Equal(t, "published", article.Status, "All articles should have 'published' status")
		}
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
		// 複数の記事を作成して公開状態にする
		slugs := CreateMultipleTestArticles(t, router, sessionToken, 3)
		PublishMultipleArticles(t, router, sessionToken, slugs)

		// リクエスト実行
		w := e2e.PerformRequest(router, http.MethodGet, "/api/articles", nil, "")

		// ステータスコードの検証
		assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")

		// レスポンスの検証
		var response dto.ArticleListResponse
		e2e.GetResponseJSON(t, w, &response)

		// 公開記事が含まれるはず
		assert.NotNil(t, response.Data, "Data field should not be nil")
		assert.NotEmpty(t, response.Data, "Articles data should not be empty")
		// 少なくとも4つの公開記事があるはず（最初に作成した1つと、このテストで作成した3つ）
		assert.GreaterOrEqual(t, len(response.Data), 4, "Should have at least 4 published articles")
		
		// すべての記事が公開状態であることを検証
		for _, article := range response.Data {
			assert.Equal(t, "published", article.Status, "All articles should have 'published' status")
		}
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

	// 他ユーザーのドラフト記事が含まれないことを確認するテスト
	t.Run("他ユーザーのドラフト記事は公開記事一覧に含まれない", func(t *testing.T) {
		// 新しいテスト環境をセットアップ
		newRouter, newCleanup := e2e.SetupTestEnvironment(t)
		defer newCleanup()

		// 最初のユーザーを作成してログイン
		_, firstUserToken, _ := e2e.CreateAndLoginTestUser(t, newRouter)

		// 最初のユーザーで公開記事を作成
		firstUserPublishedSlug := CreateTestArticle(t, newRouter, firstUserToken)
		PublishArticle(t, newRouter, firstUserToken, firstUserPublishedSlug)

		// 最初のユーザーでドラフト記事を作成
		firstUserDraftSlug := CreateTestArticle(t, newRouter, firstUserToken)

		// 別のユーザーを作成してログイン
		_, secondUserToken, _ := e2e.CreateAndLoginTestUser(t, newRouter)

		// 別のユーザーで公開記事を作成
		secondUserPublishedSlug := CreateTestArticle(t, newRouter, secondUserToken)
		PublishArticle(t, newRouter, secondUserToken, secondUserPublishedSlug)

		// 別のユーザーでドラフト記事を作成
		secondUserDraftSlug := CreateTestArticle(t, newRouter, secondUserToken)

		// リクエスト実行（認証なし）
		w := e2e.PerformRequest(newRouter, http.MethodGet, "/api/articles", nil, "")

		// ステータスコードの検証
		assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")

		// レスポンスの検証
		var response dto.ArticleListResponse
		e2e.GetResponseJSON(t, w, &response)

		// 公開記事のみが含まれるはず
		assert.NotEmpty(t, response.Data, "Articles data should not be empty")

		// ドラフト記事のSlugが含まれないことを検証
		for _, article := range response.Data {
			assert.NotEqual(t, firstUserDraftSlug, article.Slug, "First user's draft article should not be included")
			assert.NotEqual(t, secondUserDraftSlug, article.Slug, "Second user's draft article should not be included")
			// すべての記事が公開状態であることを確認
			assert.Equal(t, "published", article.Status, "All articles should have 'published' status")
		}

		// 公開記事のSlugが含まれることを検証
		slugsFound := make(map[string]bool)
		for _, article := range response.Data {
			slugsFound[article.Slug] = true
		}
		assert.True(t, slugsFound[firstUserPublishedSlug], "First user's published article should be included")
		assert.True(t, slugsFound[secondUserPublishedSlug], "Second user's published article should be included")
	})
}
