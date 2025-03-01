package auth

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/riii111/markdown-blog-api/tests/e2e"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ログアウトAPIのテスト
func TestLogoutUser(t *testing.T) {
	// テスト環境のセットアップ
	router, cleanup := e2e.SetupTestEnvironment(t)
	defer cleanup()

	// テスト用ユーザーの作成とログイン
	_, sessionToken := e2e.CreateTestUser(t, router)

	t.Run("正常系: ログアウト成功", func(t *testing.T) {
		// ログアウトリクエスト
		w := e2e.PerformRequest(router, http.MethodPost, "/api/users/logout", nil, sessionToken)

		// ステータスコードの検証
		assert.Equal(t, http.StatusNoContent, w.Code, "Expected status code 204")

		// セッションCookieの検証
		cookies := w.Result().Cookies()
		var sessionCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == "test-session" {
				sessionCookie = cookie
				break
			}
		}

		// セッションCookieが削除されているか（MaxAge < 0）
		if sessionCookie != nil {
			assert.True(t, sessionCookie.MaxAge < 0, "Session cookie should be deleted")
		}
	})

	t.Run("異常系: 認証なしでログアウト", func(t *testing.T) {
		// 認証なしでログアウトリクエスト
		w := e2e.PerformRequest(router, http.MethodPost, "/api/users/logout", nil, "")

		// ステータスコードの検証
		assert.Equal(t, http.StatusUnauthorized, w.Code, "Expected status code 401")

		// エラーレスポンスの検証
		var errorResponse map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
		require.NoError(t, err, "Failed to unmarshal error response")

		assert.Contains(t, errorResponse, "error", "Response should contain error field")
		assert.Contains(t, errorResponse["error"], "Unauthorized", "Error message should indicate unauthorized")
	})
}
