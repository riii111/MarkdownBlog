package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/riii111/markdown-blog-api/internal/handler/dto"
	"github.com/riii111/markdown-blog-api/tests/e2e"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err, "Failed to unmarshal response")

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
		var errorResponse map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
		require.NoError(t, err, "Failed to unmarshal error response")

		assert.Contains(t, errorResponse, "error", "Response should contain error field")
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
		var errorResponse map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
		require.NoError(t, err, "Failed to unmarshal error response")

		assert.Contains(t, errorResponse, "error", "Response should contain error field")
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
		var errorResponse map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
		require.NoError(t, err, "Failed to unmarshal error response")

		assert.Contains(t, errorResponse, "error", "Response should contain error field")
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
		var errorResponse map[string]string
		err := json.Unmarshal(w2.Body.Bytes(), &errorResponse)
		require.NoError(t, err, "Failed to unmarshal error response")

		assert.Contains(t, errorResponse, "error", "Response should contain error field")
	})
}
