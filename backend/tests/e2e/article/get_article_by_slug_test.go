package article

import (
	"net/http"
	"strings"
	"testing"

	"github.com/riii111/markdown-blog-api/internal/handler/dto"
	"github.com/riii111/markdown-blog-api/tests/e2e"
	"github.com/stretchr/testify/assert"
)

// 記事詳細取得APIのテスト
func TestGetArticleBySlug(t *testing.T) {
	// テスト環境のセットアップ
	router, cleanup := e2e.SetupTestEnvironment(t)
	defer cleanup()

	// テスト用ユーザーの作成とログイン
	_, sessionToken, _ := e2e.CreateAndLoginTestUser(t, router)

	// テスト環境はSetupTestEnvironmentのcleanup関数によってクリーンアップされる

	t.Run("正常系: 記事詳細取得成功", func(t *testing.T) {
		// ヘルパー関数を使用して記事作成して公開状態にする
		slug := CreateTestArticle(t, router, sessionToken)
		PublishArticle(t, router, sessionToken, slug)

		// 記事詳細取得リクエスト
		getURL := "/api/articles/" + slug
		w := e2e.PerformRequest(router, http.MethodGet, getURL, nil, "")

		// ステータスコードの検証
		assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")

		// レスポンスの検証
		var response dto.ArticleDetailResponse
		e2e.GetResponseJSON(t, w, &response)

		assert.NotEmpty(t, response.Data.ID, "Article ID should not be empty")
		assert.Equal(t, slug, response.Data.Slug, "Slug should match")
		assert.NotNil(t, response.Data.User, "User field should not be nil")
		assert.Equal(t, "published", response.Data.Status, "Article status should be 'published'")
	})

	t.Run("異常系: 存在しない記事の詳細取得", func(t *testing.T) {
		// 存在しないslugで記事詳細取得リクエスト
		nonExistentSlug := "non-existent-article-slug"
		getURL := "/api/articles/" + nonExistentSlug
		w := e2e.PerformRequest(router, http.MethodGet, getURL, nil, "")

		// エラーレスポンスの検証
		e2e.AssertErrorResponse(t, w, http.StatusNotFound, "Article not found")
	})

	t.Run("エッジケース: 特殊文字を含むslugで記事詳細取得", func(t *testing.T) {
		// 特殊文字を含むslugで記事詳細取得リクエスト
		// 注: 実際のAPIでは特殊文字を含むslugは生成されないかもしれないが、
		// エッジケースとしてテストする
		specialSlug := "special-slug-with-unusual-characters-._~"
		getURL := "/api/articles/" + specialSlug
		w := e2e.PerformRequest(router, http.MethodGet, getURL, nil, "")

		// 存在しない記事なので404が返るはず
		e2e.AssertErrorResponse(t, w, http.StatusNotFound, "Article not found")
	})

	t.Run("エッジケース: 非常に長いslugで記事詳細取得", func(t *testing.T) {
		// 非常に長いslugで記事詳細取得リクエスト
		longSlug := "very-long-slug-" + strings.Repeat("a", 200)
		getURL := "/api/articles/" + longSlug
		w := e2e.PerformRequest(router, http.MethodGet, getURL, nil, "")

		// 存在しない記事なので404が返るはず
		e2e.AssertErrorResponse(t, w, http.StatusNotFound, "Article not found")
	})

	t.Run("正常系: 複数記事作成後に各記事の詳細取得", func(t *testing.T) {
		// 複数の記事を作成して公開状態にする
		slugs := CreateMultipleTestArticles(t, router, sessionToken, 3)
		PublishMultipleArticles(t, router, sessionToken, slugs)

		// 各記事の詳細を取得して検証
		for _, slug := range slugs {
			getURL := "/api/articles/" + slug
			w := e2e.PerformRequest(router, http.MethodGet, getURL, nil, "")

			// ステータスコードの検証
			assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")

			// レスポンスの検証
			var response dto.ArticleDetailResponse
			e2e.GetResponseJSON(t, w, &response)

			assert.NotEmpty(t, response.Data.ID, "Article ID should not be empty")
			assert.Equal(t, slug, response.Data.Slug, "Slug should match")
		}
	})

	// 他ユーザーのドラフト記事は取得できないことを検証するテスト
	t.Run("他ユーザーのドラフト記事は取得できない", func(t *testing.T) {
		// 新しいテスト環境をセットアップ
		newRouter, newCleanup := e2e.SetupTestEnvironment(t)
		defer newCleanup()

		// 最初のユーザーを作成してログイン
		_, firstUserToken, _ := e2e.CreateAndLoginTestUser(t, newRouter)

		// 別のユーザーを作成してログイン
		_, secondUserToken, _ := e2e.CreateAndLoginTestUser(t, newRouter)

		// 別のユーザーでドラフト記事を作成
		secondUserDraftSlug := CreateTestArticle(t, newRouter, secondUserToken)

		// 最初のユーザーとして、別ユーザーのドラフト記事へアクセス
		getURL := "/api/articles/" + secondUserDraftSlug
		w := e2e.PerformRequest(newRouter, http.MethodGet, getURL, nil, firstUserToken)

		// ドラフト記事は他ユーザーがアクセスできないので404が返るはず
		assert.Equal(t, http.StatusNotFound, w.Code, "Expected status code 404 when accessing another user's draft article")

		// 認証なしでアクセスしても404になることを確認
		w = e2e.PerformRequest(newRouter, http.MethodGet, getURL, nil, "")
		assert.Equal(t, http.StatusNotFound, w.Code, "Expected status code 404 when accessing a draft article without authentication")

		// ドラフト記事作成者は自分のドラフト記事にアクセスできることを確認
		w = e2e.PerformRequest(newRouter, http.MethodGet, getURL, nil, secondUserToken)
		assert.Equal(t, http.StatusOK, w.Code, "Owner should be able to access their own draft article")

		// レスポンスの検証
		var response dto.ArticleDetailResponse
		e2e.GetResponseJSON(t, w, &response)

		assert.NotEmpty(t, response.Data.ID, "Article ID should not be empty")
		assert.Equal(t, secondUserDraftSlug, response.Data.Slug, "Slug should match")
		assert.Equal(t, "draft", response.Data.Status, "Article status should be 'draft'")
	})

	// 公開記事は誰でも取得できることを検証するテスト
	t.Run("公開記事は誰でも取得できる", func(t *testing.T) {
		// 新しいテスト環境をセットアップ
		newRouter, newCleanup := e2e.SetupTestEnvironment(t)
		defer newCleanup()

		// 最初のユーザーを作成してログイン
		_, firstUserToken, _ := e2e.CreateAndLoginTestUser(t, newRouter)

		// 別のユーザーを作成してログイン
		_, secondUserToken, _ := e2e.CreateAndLoginTestUser(t, newRouter)

		// 最初のユーザーで公開記事を作成
		firstUserPublishedSlug := CreateTestArticle(t, newRouter, firstUserToken)
		PublishArticle(t, newRouter, firstUserToken, firstUserPublishedSlug)

		// 認証なしでアクセスしても200になることを確認
		getURL := "/api/articles/" + firstUserPublishedSlug
		w := e2e.PerformRequest(newRouter, http.MethodGet, getURL, nil, "")
		assert.Equal(t, http.StatusOK, w.Code, "Anyone should be able to access a published article without authentication")

		// 別ユーザーとしてアクセスしても200になることを確認
		w = e2e.PerformRequest(newRouter, http.MethodGet, getURL, nil, secondUserToken)
		assert.Equal(t, http.StatusOK, w.Code, "Any authenticated user should be able to access another user's published article")

		// レスポンスの検証
		var response dto.ArticleDetailResponse
		e2e.GetResponseJSON(t, w, &response)

		assert.NotEmpty(t, response.Data.ID, "Article ID should not be empty")
		assert.Equal(t, firstUserPublishedSlug, response.Data.Slug, "Slug should match")
		assert.Equal(t, "published", response.Data.Status, "Article status should be 'published'")
	})
}
