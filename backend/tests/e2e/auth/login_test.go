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
		e2e.GetResponseJSON(t, w, &response)

		assert.NotEmpty(t, response.ID, "User ID should not be empty")
		assert.Equal(t, testUser.DisplayName, response.DisplayName, "Display name should match")

		// セッションCookieの検証
		cookies := w.Result().Cookies()
		var sessionCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == "test-session" {
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

		// エラーレスポンスの検証
		e2e.AssertErrorResponse(t, w, http.StatusUnauthorized, "Invalid credentials")
	})

	t.Run("異常系: 存在しないユーザー", func(t *testing.T) {
		// 存在しないユーザーでログインリクエスト
		loginReq := dto.LoginRequest{
			Email:    "nonexistent@example.com",
			Password: "password123",
		}

		// リクエスト実行
		w := e2e.PerformRequest(router, http.MethodPost, "/api/users/login", loginReq, "")

		// エラーレスポンスの検証
		e2e.AssertErrorResponse(t, w, http.StatusUnauthorized, "Invalid credentials")
	})

	t.Run("異常系: 空のリクエストボディ", func(t *testing.T) {
		// 空のリクエストボディでログイン試行
		w := e2e.PerformRequest(router, http.MethodPost, "/api/users/login", nil, "")

		// エラーレスポンスの検証
		e2e.AssertErrorResponse(t, w, http.StatusBadRequest, "")
	})

	t.Run("異常系: 不正なJSONフォーマット", func(t *testing.T) {
		// 不正なJSONフォーマットでリクエスト
		invalidJSON := `{"email": "test@example.com", "password": missing_quotes}`

		// 直接HTTPリクエストを作成して実行
		req := httptest.NewRequest(http.MethodPost, "/api/users/login", strings.NewReader(invalidJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// エラーレスポンスの検証
		e2e.AssertErrorResponse(t, w, http.StatusBadRequest, "")
	})

	t.Run("異常系: メールアドレスのみ提供（パスワード欠落）", func(t *testing.T) {
		// パスワードが欠落したリクエスト
		incompleteReq := map[string]string{
			"email": testUser.Email,
		}

		// リクエスト実行
		w := e2e.PerformRequest(router, http.MethodPost, "/api/users/login", incompleteReq, "")

		// エラーレスポンスの検証
		e2e.AssertErrorResponse(t, w, http.StatusBadRequest, "")
	})

	t.Run("異常系: パスワードのみ提供（メールアドレス欠落）", func(t *testing.T) {
		// メールアドレスが欠落したリクエスト
		incompleteReq := map[string]string{
			"password": testUser.Password,
		}

		// リクエスト実行
		w := e2e.PerformRequest(router, http.MethodPost, "/api/users/login", incompleteReq, "")

		// エラーレスポンスの検証
		e2e.AssertErrorResponse(t, w, http.StatusBadRequest, "")
	})
}
