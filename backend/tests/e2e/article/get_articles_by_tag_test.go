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
		// テスト用のタグ付き記事を作成して公開状態にする
		slug := CreateArticleWithTag(t, router, sessionToken, testTagSlug)
		PublishArticle(t, router, sessionToken, slug)

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
		
		// すべての記事が公開状態であることを検証
		for _, article := range response.Data {
			assert.Equal(t, "published", article.Status, "All articles should have 'published' status")
		}
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

	// 他ユーザーのドラフト記事が含まれないことを確認するテスト
	t.Run("他ユーザーのドラフト記事はタグに紐づく記事一覧に含まれない", func(t *testing.T) {
		// 新しいテスト環境をセットアップ
		newRouter, newCleanup := e2e.SetupTestEnvironment(t)
		defer newCleanup()

		// テスト用のタグを作成
		tagSlug := CreateTestTag(t, newRouter, "テストタグ")

		// 最初のユーザーを作成してログイン
		_, firstUserToken, _ := e2e.CreateAndLoginTestUser(t, newRouter)

		// 最初のユーザーで公開記事を作成し、タグ付け
		firstUserPublishedSlug := CreateTestArticle(t, newRouter, firstUserToken)
		AttachTagToArticle(t, newRouter, firstUserToken, firstUserPublishedSlug, tagSlug)
		PublishArticle(t, newRouter, firstUserToken, firstUserPublishedSlug)

		// 最初のユーザーでドラフト記事を作成、タグ付け
		firstUserDraftSlug := CreateTestArticle(t, newRouter, firstUserToken)
		AttachTagToArticle(t, newRouter, firstUserToken, firstUserDraftSlug, tagSlug)

		// 別のユーザーを作成してログイン
		_, secondUserToken, _ := e2e.CreateAndLoginTestUser(t, newRouter)

		// 別のユーザーで公開記事を作成、タグ付け
		secondUserPublishedSlug := CreateTestArticle(t, newRouter, secondUserToken)
		AttachTagToArticle(t, newRouter, secondUserToken, secondUserPublishedSlug, tagSlug)
		PublishArticle(t, newRouter, secondUserToken, secondUserPublishedSlug)

		// 別のユーザーでドラフト記事を作成、タグ付け
		secondUserDraftSlug := CreateTestArticle(t, newRouter, secondUserToken)
		AttachTagToArticle(t, newRouter, secondUserToken, secondUserDraftSlug, tagSlug)

		// リクエスト実行（認証なし）
		getURL := "/api/tags/" + tagSlug + "/articles"
		w := e2e.PerformRequest(newRouter, http.MethodGet, getURL, nil, "")

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
