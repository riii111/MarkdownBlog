package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/riii111/markdown-blog-api/internal/handler/dto"
	"github.com/riii111/markdown-blog-api/tests/e2e"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ログインAPIのテスト
func TestLoginUser(t *testing.T) {
	// テスト環境のセットアップ
	router, cleanup := e2e.SetupTestEnvironment(t)
	defer cleanup()

	// テスト用ユーザーの作成
	testUser := dto.RegisterUserRequest{
		Email:       fmt.Sprintf("login-test-%s@example.com", uuid.New().String()),
		Password:    "password123",
		DisplayName: "Login Test User",
	}

	// ユーザー登録
	w := e2e.PerformRequest(router, http.MethodPost, "/api/users/register", testUser, "")
	require.Equal(t, http.StatusCreated, w.Code, "User registration should succeed")

	t.Run("正常系: 有効な認証情報でログイン", func(t *testing.T) {
		// ログインリクエスト
		loginReq := dto.LoginRequest{
			Email:    testUser.Email,
			Password: testUser.Password,
		}

		// リクエスト実行
		w := e2e.PerformRequest(router, http.MethodPost, "/api/users/login", loginReq, "")

		// ステータスコードの検証
		assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")

		// レスポンスの検証
		var response dto.LoginResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err, "Failed to unmarshal response")

		assert.NotEmpty(t, response.ID, "User ID should not be empty")
		assert.Equal(t, testUser.DisplayName, response.DisplayName, "Display name should match")

		// セッションCookieの検証
		cookies := w.Result().Cookies()
		var sessionCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == os.Getenv("SESSION_NAME") {
				sessionCookie = cookie
				break
			}
		}
		assert.NotNil(t, sessionCookie, "Session cookie should be set")
	})

	t.Run("異常系: 無効なパスワード", func(t *testing.T) {
		// 無効なパスワードでログインリクエスト
		loginReq := dto.LoginRequest{
			Email:    testUser.Email,
			Password: "wrongpassword",
		}

		// リクエスト実行
		w := e2e.PerformRequest(router, http.MethodPost, "/api/users/login", loginReq, "")

		// ステータスコードの検証
		assert.Equal(t, http.StatusUnauthorized, w.Code, "Expected status code 401")

		// エラーレスポンスの検証
		var errorResponse map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
		require.NoError(t, err, "Failed to unmarshal error response")

		assert.Contains(t, errorResponse, "error", "Response should contain error field")
		assert.Contains(t, errorResponse["error"], "Invalid credentials", "Error message should indicate invalid credentials")
	})

	t.Run("異常系: 存在しないユーザー", func(t *testing.T) {
		// 存在しないユーザーでログインリクエスト
		loginReq := dto.LoginRequest{
			Email:    "nonexistent@example.com",
			Password: "password123",
		}

		// リクエスト実行
		w := e2e.PerformRequest(router, http.MethodPost, "/api/users/login", loginReq, "")

		// ステータスコードの検証
		assert.Equal(t, http.StatusUnauthorized, w.Code, "Expected status code 401")

		// エラーレスポンスの検証
		var errorResponse map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
		require.NoError(t, err, "Failed to unmarshal error response")

		assert.Contains(t, errorResponse, "error", "Response should contain error field")
		assert.Contains(t, errorResponse["error"], "Invalid credentials", "Error message should indicate invalid credentials")
	})
}
