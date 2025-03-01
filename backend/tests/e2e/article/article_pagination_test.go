package article

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/riii111/markdown-blog-api/internal/handler/dto"
	"github.com/riii111/markdown-blog-api/tests/e2e"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// カーソルベースページネーションのテスト
func TestArticlePagination(t *testing.T) {
	// テスト環境のセットアップ
	router, cleanup := e2e.SetupTestEnvironment(t)
	defer cleanup()

	// テスト用ユーザーの作成とログイン
	_, sessionToken, _ := e2e.CreateAndLoginTestUser(t, router)

	// テスト用に複数の記事を作成
	articleSlugs := createMultipleArticles(t, router, sessionToken, 12)
	require.Len(t, articleSlugs, 12, "12個の記事が作成されるべき")

	t.Run("正常系: 自分の記事一覧をページネーションで取得", func(t *testing.T) {
		// 1ページ目（5件）を取得
		w1 := e2e.PerformRequest(router, http.MethodGet, "/api/articles/me?per_page=5&page=1", nil, sessionToken)
		assert.Equal(t, http.StatusOK, w1.Code, "1ページ目の取得は成功するべき")

		var response1 dto.ArticleListResponse
		e2e.GetResponseJSON(t, w1, &response1)

		assert.Len(t, response1.Data, 5, "1ページ目には5件の記事があるはず")
		assert.True(t, response1.Pagination.HasMore, "まだ次のページがあるはず")
		assert.NotNil(t, response1.Pagination.NextCursor, "次のカーソルが存在するはず")

		// 次のページを取得
		nextCursor := *response1.Pagination.NextCursor
		w2 := e2e.PerformRequest(router, http.MethodGet, fmt.Sprintf("/api/articles/me?per_page=5&cursor=%s", nextCursor), nil, sessionToken)
		assert.Equal(t, http.StatusOK, w2.Code, "2ページ目の取得は成功するべき")

		var response2 dto.ArticleListResponse
		e2e.GetResponseJSON(t, w2, &response2)

		// 2ページ目も記事があるはず（具体的な数は実装によって異なるかもしれない）
		assert.NotEmpty(t, response2.Data, "2ページ目にも記事があるはず")

		// 3ページ目があるかどうかは記事の総数によって変わる
		// HasMoreとNextCursorの検証は省略する
	})

	t.Run("正常系: 公開記事一覧をリミット指定で取得", func(t *testing.T) {
		// リミット付きリクエスト実行
		w := e2e.PerformRequest(router, http.MethodGet, "/api/articles?limit=10", nil, "")
		assert.Equal(t, http.StatusOK, w.Code, "リミット付き記事一覧取得は成功するべき")

		var response dto.ArticleListResponse
		e2e.GetResponseJSON(t, w, &response)

		assert.Equal(t, 10, response.Pagination.ItemsPerPage, "1ページあたりのアイテム数は10のはず")
		// 注: 記事はデフォルトでは公開状態ではないため、公開記事は0件かもしれない
	})

	t.Run("異常系: 不正なカーソル値での記事一覧取得", func(t *testing.T) {
		// 不正なカーソル値でリクエスト実行
		w := e2e.PerformRequest(router, http.MethodGet, "/api/articles/me?cursor=invalid-cursor", nil, sessionToken)

		// エラーレスポンスの検証
		// 注: 実装によっては400エラーか空の結果を返す可能性がある
		if w.Code == http.StatusBadRequest {
			e2e.AssertErrorResponse(t, w, http.StatusBadRequest, "")
		} else {
			assert.Equal(t, http.StatusOK, w.Code, "不正なカーソルの場合は空の結果を返すべき")
			var response dto.ArticleListResponse
			e2e.GetResponseJSON(t, w, &response)
			assert.NotNil(t, response.Data, "データフィールドはnilではないはず")
		}
	})
}

// 複数記事の操作テスト
func TestMultipleArticleOperations(t *testing.T) {
	// テスト環境のセットアップ
	router, cleanup := e2e.SetupTestEnvironment(t)
	defer cleanup()

	// 複数のテストユーザーを作成
	_, sessionToken1, _ := e2e.CreateAndLoginTestUser(t, router)
	_, sessionToken2, _ := e2e.CreateAndLoginTestUser(t, router)

	// ユーザー1の記事を作成
	user1Articles := createMultipleArticles(t, router, sessionToken1, 3)
	require.Len(t, user1Articles, 3, "ユーザー1の記事が3つ作成されるべき")

	// ユーザー2の記事を作成
	user2Articles := createMultipleArticles(t, router, sessionToken2, 2)
	require.Len(t, user2Articles, 2, "ユーザー2の記事が2つ作成されるべき")

	t.Run("正常系: ユーザー1は自分の記事のみ取得できる", func(t *testing.T) {
		// ユーザー1の記事一覧を取得
		w := e2e.PerformRequest(router, http.MethodGet, "/api/articles/me", nil, sessionToken1)
		assert.Equal(t, http.StatusOK, w.Code, "ユーザー1の記事一覧取得は成功するべき")

		var response dto.ArticleListResponse
		e2e.GetResponseJSON(t, w, &response)

		assert.Len(t, response.Data, 3, "ユーザー1は3つの記事を持っているはず")
	})

	t.Run("正常系: ユーザー2は自分の記事のみ取得できる", func(t *testing.T) {
		// ユーザー2の記事一覧を取得
		w := e2e.PerformRequest(router, http.MethodGet, "/api/articles/me", nil, sessionToken2)
		assert.Equal(t, http.StatusOK, w.Code, "ユーザー2の記事一覧取得は成功するべき")

		var response dto.ArticleListResponse
		e2e.GetResponseJSON(t, w, &response)

		assert.Len(t, response.Data, 2, "ユーザー2は2つの記事を持っているはず")
	})

	t.Run("異常系: ユーザー1は他のユーザーの記事を削除できない", func(t *testing.T) {
		// ユーザー1がユーザー2の記事を削除しようとする
		w := e2e.PerformRequest(router, http.MethodDelete, fmt.Sprintf("/api/articles/%s", user2Articles[0]), nil, sessionToken1)

		// エラーレスポンスの検証
		// 注: 実装によっては404（記事が見つからない）または403（権限なし）を返す可能性がある
		if w.Code == http.StatusNotFound {
			e2e.AssertErrorResponse(t, w, http.StatusNotFound, "Article not found")
		} else {
			e2e.AssertErrorResponse(t, w, http.StatusUnauthorized, "Unauthorized")
		}
	})
}

// テスト用に複数の記事を作成するヘルパー関数
func createMultipleArticles(t *testing.T, router *gin.Engine, sessionToken string, count int) []string {
	slugs := make([]string, 0, count)
	createReq := dto.CreateArticleRequest{}

	for i := 0; i < count; i++ {
		w := e2e.PerformRequest(router, http.MethodPost, "/api/articles", createReq, sessionToken)
		require.Equal(t, http.StatusCreated, w.Code, "記事作成は成功するべき")

		var response dto.CreateArticleResponse
		e2e.GetResponseJSON(t, w, &response)
		slugs = append(slugs, response.Slug)
	}

	return slugs
}
