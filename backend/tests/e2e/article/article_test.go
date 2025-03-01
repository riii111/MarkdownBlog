package article

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/riii111/markdown-blog-api/internal/handler/dto"
	"github.com/riii111/markdown-blog-api/tests/e2e"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 記事作成APIのテスト
func TestCreateArticle(t *testing.T) {
	// テスト環境のセットアップ
	router, cleanup := e2e.SetupTestEnvironment(t)
	defer cleanup()

	// テスト用ユーザーの作成とログイン
	_, sessionToken, _ := e2e.CreateAndLoginTestUser(t, router)

	t.Run("正常系: 記事の作成成功", func(t *testing.T) {
		// 記事作成リクエスト
		createReq := dto.CreateArticleRequest{}

		// リクエスト実行
		w := e2e.PerformRequest(router, http.MethodPost, "/api/articles", createReq, sessionToken)

		// ステータスコードの検証
		assert.Equal(t, http.StatusCreated, w.Code, "記事作成は成功するべき")

		// レスポンスの検証
		var response dto.CreateArticleResponse
		e2e.GetResponseJSON(t, w, &response)

		assert.NotEmpty(t, response.Slug, "記事のスラグは空ではないはず")
	})

	t.Run("異常系: 認証なしで記事作成", func(t *testing.T) {
		// 認証なしで記事作成リクエスト
		createReq := dto.CreateArticleRequest{}

		// リクエスト実行
		w := e2e.PerformRequest(router, http.MethodPost, "/api/articles", createReq, "")

		// エラーレスポンスの検証
		e2e.AssertErrorResponse(t, w, http.StatusUnauthorized, "Unauthorized")
	})
}

// 記事削除APIのテスト
func TestDeleteArticle(t *testing.T) {
	// テスト環境のセットアップ
	router, cleanup := e2e.SetupTestEnvironment(t)
	defer cleanup()

	// テスト用ユーザーの作成とログイン
	_, sessionToken, _ := e2e.CreateAndLoginTestUser(t, router)

	// テスト用記事の作成
	createReq := dto.CreateArticleRequest{}
	w := e2e.PerformRequest(router, http.MethodPost, "/api/articles", createReq, sessionToken)
	require.Equal(t, http.StatusCreated, w.Code, "記事作成は成功するべき")

	var createResp dto.CreateArticleResponse
	e2e.GetResponseJSON(t, w, &createResp)
	articleSlug := createResp.Slug

	t.Run("正常系: 記事の削除成功", func(t *testing.T) {
		// 記事削除リクエスト
		w := e2e.PerformRequest(router, http.MethodDelete, fmt.Sprintf("/api/articles/%s", articleSlug), nil, sessionToken)

		// ステータスコードの検証
		assert.Equal(t, http.StatusNoContent, w.Code, "記事削除は成功するべき")
	})

	t.Run("異常系: 存在しない記事の削除", func(t *testing.T) {
		// 存在しない記事のスラグ
		nonExistentSlug := "non-existent-slug"

		// リクエスト実行
		w := e2e.PerformRequest(router, http.MethodDelete, fmt.Sprintf("/api/articles/%s", nonExistentSlug), nil, sessionToken)

		// エラーレスポンスの検証
		e2e.AssertErrorResponse(t, w, http.StatusNotFound, "Article not found")
	})

	t.Run("異常系: 認証なしで記事削除", func(t *testing.T) {
		// 新しい記事を作成
		createW := e2e.PerformRequest(router, http.MethodPost, "/api/articles", createReq, sessionToken)
		var newArticleResp dto.CreateArticleResponse
		e2e.GetResponseJSON(t, createW, &newArticleResp)

		// 認証なしで記事削除リクエスト
		w := e2e.PerformRequest(router, http.MethodDelete, fmt.Sprintf("/api/articles/%s", newArticleResp.Slug), nil, "")

		// エラーレスポンスの検証
		e2e.AssertErrorResponse(t, w, http.StatusUnauthorized, "Unauthorized")
	})
}

// 記事一覧取得APIのテスト
func TestGetArticles(t *testing.T) {
	// テスト環境のセットアップ
	router, cleanup := e2e.SetupTestEnvironment(t)
	defer cleanup()

	t.Run("正常系: 記事一覧の取得", func(t *testing.T) {
		// リクエスト実行
		w := e2e.PerformRequest(router, http.MethodGet, "/api/articles", nil, "")

		// ステータスコードの検証
		assert.Equal(t, http.StatusOK, w.Code, "記事一覧取得は成功するべき")

		// レスポンスの検証
		var response dto.ArticleListResponse
		e2e.GetResponseJSON(t, w, &response)

		// 初期状態では記事がないかもしれないが、レスポンス構造は正しいはず
		assert.NotNil(t, response.Data, "データフィールドはnilではないはず")
		assert.NotNil(t, response.Pagination, "ページネーションフィールドはnilではないはず")
	})

	t.Run("正常系: リミット指定での記事一覧取得", func(t *testing.T) {
		// リミット付きリクエスト実行
		w := e2e.PerformRequest(router, http.MethodGet, "/api/articles?limit=5", nil, "")

		// ステータスコードの検証
		assert.Equal(t, http.StatusOK, w.Code, "リミット付き記事一覧取得は成功するべき")

		// レスポンスの検証
		var response dto.ArticleListResponse
		e2e.GetResponseJSON(t, w, &response)

		assert.Equal(t, 5, response.Pagination.ItemsPerPage, "1ページあたりのアイテム数は5のはず")
	})
}

// 記事詳細取得APIのテスト
func TestGetArticleBySlug(t *testing.T) {
	// テスト環境のセットアップ
	router, cleanup := e2e.SetupTestEnvironment(t)
	defer cleanup()

	// テスト用ユーザーの作成とログイン
	_, sessionToken, _ := e2e.CreateAndLoginTestUser(t, router)

	// テスト用記事の作成
	createReq := dto.CreateArticleRequest{}
	w := e2e.PerformRequest(router, http.MethodPost, "/api/articles", createReq, sessionToken)
	require.Equal(t, http.StatusCreated, w.Code, "記事作成は成功するべき")

	var createResp dto.CreateArticleResponse
	e2e.GetResponseJSON(t, w, &createResp)
	articleSlug := createResp.Slug

	t.Run("正常系: 記事詳細の取得", func(t *testing.T) {
		// 記事詳細取得リクエスト
		w := e2e.PerformRequest(router, http.MethodGet, fmt.Sprintf("/api/articles/%s", articleSlug), nil, "")

		// ステータスコードの検証
		assert.Equal(t, http.StatusOK, w.Code, "記事詳細取得は成功するべき")

		// レスポンスの検証
		var response dto.ArticleDetailResponse
		e2e.GetResponseJSON(t, w, &response)

		assert.Equal(t, articleSlug, response.Data.Slug, "取得した記事のスラグは一致するはず")
	})

	t.Run("異常系: 存在しない記事の詳細取得", func(t *testing.T) {
		// 存在しない記事のスラグ
		nonExistentSlug := "non-existent-slug"

		// リクエスト実行
		w := e2e.PerformRequest(router, http.MethodGet, fmt.Sprintf("/api/articles/%s", nonExistentSlug), nil, "")

		// エラーレスポンスの検証
		e2e.AssertErrorResponse(t, w, http.StatusNotFound, "Article not found")
	})
}

// 自分の記事一覧取得APIのテスト
func TestGetMeArticles(t *testing.T) {
	// テスト環境のセットアップ
	router, cleanup := e2e.SetupTestEnvironment(t)
	defer cleanup()

	// テスト用ユーザーの作成とログイン
	_, sessionToken, _ := e2e.CreateAndLoginTestUser(t, router)

	// テスト用記事の作成
	createReq := dto.CreateArticleRequest{}
	w := e2e.PerformRequest(router, http.MethodPost, "/api/articles", createReq, sessionToken)
	require.Equal(t, http.StatusCreated, w.Code, "記事作成は成功するべき")

	t.Run("正常系: 自分の記事一覧の取得", func(t *testing.T) {
		// 自分の記事一覧取得リクエスト
		w := e2e.PerformRequest(router, http.MethodGet, "/api/articles/me", nil, sessionToken)

		// ステータスコードの検証
		assert.Equal(t, http.StatusOK, w.Code, "自分の記事一覧取得は成功するべき")

		// レスポンスの検証
		var response dto.ArticleListResponse
		e2e.GetResponseJSON(t, w, &response)

		assert.NotNil(t, response.Data, "データフィールドはnilではないはず")
		assert.GreaterOrEqual(t, len(response.Data), 1, "少なくとも1つの記事があるはず")
	})

	t.Run("異常系: 認証なしで自分の記事一覧取得", func(t *testing.T) {
		// 認証なしで自分の記事一覧取得リクエスト
		w := e2e.PerformRequest(router, http.MethodGet, "/api/articles/me", nil, "")

		// エラーレスポンスの検証
		e2e.AssertErrorResponse(t, w, http.StatusUnauthorized, "Unauthorized")
	})

	t.Run("正常系: ページネーションパラメータ付きで自分の記事一覧取得", func(t *testing.T) {
		// ページネーションパラメータ付きリクエスト実行
		w := e2e.PerformRequest(router, http.MethodGet, "/api/articles/me?per_page=5&page=1", nil, sessionToken)

		// ステータスコードの検証
		assert.Equal(t, http.StatusOK, w.Code, "ページネーション付き自分の記事一覧取得は成功するべき")

		// レスポンスの検証
		var response dto.ArticleListResponse
		e2e.GetResponseJSON(t, w, &response)

		assert.Equal(t, 5, response.Pagination.ItemsPerPage, "1ページあたりのアイテム数は5のはず")
	})
}

// タグに紐づく記事一覧取得APIのテスト
// TODO: このテストはタグ機能の実装が必要
func TestGetArticlesByTag(t *testing.T) {
	// テスト環境のセットアップ
	router, cleanup := e2e.SetupTestEnvironment(t)
	defer cleanup()

	t.Run("異常系: 存在しないタグの記事一覧取得", func(t *testing.T) {
		// 存在しないタグのスラグ
		nonExistentTagSlug := fmt.Sprintf("non-existent-tag-%s", uuid.New().String())

		// リクエスト実行
		w := e2e.PerformRequest(router, http.MethodGet, fmt.Sprintf("/api/tags/%s/articles", nonExistentTagSlug), nil, "")

		// ステータスコードの検証
		assert.Equal(t, http.StatusNotFound, w.Code, "存在しないタグの記事一覧取得は404を返すべき")

		// 注: レスポンスの形式はAPIの実装によって異なる可能性があるため、
		// 単純にステータスコードだけを検証する
	})
}
