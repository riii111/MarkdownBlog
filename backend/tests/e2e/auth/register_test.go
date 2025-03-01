package auth

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/riii111/markdown-blog-api/internal/handler/dto"
	"github.com/riii111/markdown-blog-api/tests/e2e"
	"github.com/stretchr/testify/assert"
)

// ユーザー登録APIのテスト
func TestRegisterUser(t *testing.T) {
	// テスト環境のセットアップ
	router, cleanup := e2e.SetupTestEnvironment(t)
	defer cleanup()

	t.Run("正常系: 有効なデータでユーザー登録", func(t *testing.T) {
		// テスト用ユーザーデータ
		testUser := dto.RegisterUserRequest{
			Email:       fmt.Sprintf("test-%s@example.com", uuid.New().String()),
			Password:    "password123",
			DisplayName: "Test User",
		}

		// リクエスト実行
		w := e2e.PerformRequest(router, http.MethodPost, "/api/users/register", testUser, "")

		// ステータスコードの検証
		assert.Equal(t, http.StatusCreated, w.Code, "Expected status code 201")

		// レスポンスの検証
		var response dto.RegisterUserResponse
		e2e.GetResponseJSON(t, w, &response)

		assert.NotEmpty(t, response.ID, "User ID should not be empty")
		assert.Equal(t, testUser.DisplayName, response.DisplayName, "Display name should match")
		assert.NotZero(t, response.CreatedAt, "Created at should not be zero")
	})

	t.Run("異常系: バリデーションエラー（メールアドレス不正）", func(t *testing.T) {
		// 不正なメールアドレスを持つユーザーデータ
		invalidUser := dto.RegisterUserRequest{
			Email:       "invalid-email",
			Password:    "password123",
			DisplayName: "Test User",
		}

		// リクエスト実行
		w := e2e.PerformRequest(router, http.MethodPost, "/api/users/register", invalidUser, "")

		// ステータスコードの検証
		assert.Equal(t, http.StatusBadRequest, w.Code, "Expected status code 400")

		// エラーレスポンスの検証
		e2e.AssertErrorResponse(t, w, http.StatusBadRequest, "")
	})

	t.Run("異常系: バリデーションエラー（パスワード短すぎ）", func(t *testing.T) {
		// パスワードが短すぎるユーザーデータ
		invalidUser := dto.RegisterUserRequest{
			Email:       fmt.Sprintf("test-%s@example.com", uuid.New().String()),
			Password:    "pass",
			DisplayName: "Test User",
		}

		// リクエスト実行
		w := e2e.PerformRequest(router, http.MethodPost, "/api/users/register", invalidUser, "")

		// ステータスコードの検証
		assert.Equal(t, http.StatusBadRequest, w.Code, "Expected status code 400")

		// エラーレスポンスの検証
		e2e.AssertErrorResponse(t, w, http.StatusBadRequest, "")
	})

	t.Run("異常系: バリデーションエラー（パスワードが長すぎる）", func(t *testing.T) {
		// パスワードが長すぎるユーザーデータ
		invalidUser := dto.RegisterUserRequest{
			Email:       fmt.Sprintf("test-%s@example.com", uuid.New().String()),
			Password:    strings.Repeat("a", 100), // 100文字のパスワード
			DisplayName: "Test User",
		}

		// リクエスト実行
		w := e2e.PerformRequest(router, http.MethodPost, "/api/users/register", invalidUser, "")

		// ステータスコードの検証
		assert.Equal(t, http.StatusBadRequest, w.Code, "Expected status code 400")

		// エラーレスポンスの検証
		e2e.AssertErrorResponse(t, w, http.StatusBadRequest, "")
	})

	t.Run("異常系: バリデーションエラー（表示名なし）", func(t *testing.T) {
		// 表示名がないユーザーデータ
		invalidUser := dto.RegisterUserRequest{
			Email:       fmt.Sprintf("test-%s@example.com", uuid.New().String()),
			Password:    "password123",
			DisplayName: "",
		}

		// リクエスト実行
		w := e2e.PerformRequest(router, http.MethodPost, "/api/users/register", invalidUser, "")

		// ステータスコードの検証
		assert.Equal(t, http.StatusBadRequest, w.Code, "Expected status code 400")

		// エラーレスポンスの検証
		e2e.AssertErrorResponse(t, w, http.StatusBadRequest, "")
	})

	t.Run("異常系: バリデーションエラー（表示名が長すぎる）", func(t *testing.T) {
		// 表示名が長すぎるユーザーデータ
		invalidUser := dto.RegisterUserRequest{
			Email:       fmt.Sprintf("test-%s@example.com", uuid.New().String()),
			Password:    "password123",
			DisplayName: strings.Repeat("a", 150), // 150文字の表示名
		}

		// リクエスト実行
		w := e2e.PerformRequest(router, http.MethodPost, "/api/users/register", invalidUser, "")

		// ステータスコードの検証
		assert.Equal(t, http.StatusBadRequest, w.Code, "Expected status code 400")

		// エラーレスポンスの検証
		e2e.AssertErrorResponse(t, w, http.StatusBadRequest, "")
	})

	t.Run("異常系: 重複メールアドレス", func(t *testing.T) {
		// 最初のユーザー登録
		testUser := dto.RegisterUserRequest{
			Email:       fmt.Sprintf("duplicate-%s@example.com", uuid.New().String()),
			Password:    "password123",
			DisplayName: "Test User",
		}

		// 1回目のリクエスト実行
		w1 := e2e.PerformRequest(router, http.MethodPost, "/api/users/register", testUser, "")
		assert.Equal(t, http.StatusCreated, w1.Code, "First registration should succeed")

		// 同じメールアドレスで2回目のリクエスト実行
		w2 := e2e.PerformRequest(router, http.MethodPost, "/api/users/register", testUser, "")

		// ステータスコードの検証
		assert.Equal(t, http.StatusBadRequest, w2.Code, "Expected status code 400 for duplicate email")

		// エラーレスポンスの検証
		e2e.AssertErrorResponse(t, w2, http.StatusBadRequest, "")
	})

	t.Run("異常系: 空のリクエストボディ", func(t *testing.T) {
		// 空のリクエストボディで登録試行
		w := e2e.PerformRequest(router, http.MethodPost, "/api/users/register", nil, "")

		// ステータスコードの検証
		assert.Equal(t, http.StatusBadRequest, w.Code, "Expected status code 400")

		// エラーレスポンスの検証
		e2e.AssertErrorResponse(t, w, http.StatusBadRequest, "")
	})

	t.Run("異常系: 不正なJSONフォーマット", func(t *testing.T) {
		// 不正なJSONフォーマットでリクエスト
		invalidJSON := `{"email": "test@example.com", "password": missing_quotes, "display_name": "Test User"}`

		// カスタムリクエスト実行（不正なJSONの場合は直接HTTPリクエストを作成）
		req := httptest.NewRequest(http.MethodPost, "/api/users/register", strings.NewReader(invalidJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// ステータスコードの検証
		assert.Equal(t, http.StatusBadRequest, w.Code, "Expected status code 400")

		// エラーレスポンスの検証
		e2e.AssertErrorResponse(t, w, http.StatusBadRequest, "")
	})
}
