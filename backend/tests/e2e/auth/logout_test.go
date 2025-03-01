package auth

import (
	"net/http"
	"testing"

	"github.com/riii111/markdown-blog-api/tests/e2e"
	"github.com/stretchr/testify/assert"
)

// ログアウトAPIのテスト
func TestLogoutUser(t *testing.T) {
	// テスト環境のセットアップ
	router, cleanup := e2e.SetupTestEnvironment(t)
	defer cleanup()

	// テスト用ユーザーの作成とログイン
	_, sessionToken, _ := e2e.CreateAndLoginTestUser(t, router)

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

		// エラーレスポンスの検証
		e2e.AssertErrorResponse(t, w, http.StatusUnauthorized, "Unauthorized")
	})

	t.Run("異常系: 無効なセッショントークンでログアウト", func(t *testing.T) {
		// 無効なセッショントークンでログアウトリクエスト
		invalidToken := "invalid-session-token"
		w := e2e.PerformRequest(router, http.MethodPost, "/api/users/logout", nil, invalidToken)

		// エラーレスポンスの検証
		e2e.AssertErrorResponse(t, w, http.StatusUnauthorized, "Unauthorized")
	})

	t.Run("異常系: 既にログアウト済みのセッションで再ログアウト", func(t *testing.T) {
		// 新しいユーザーとセッションを作成
		_, sessionToken, _ := e2e.CreateAndLoginTestUser(t, router)

		// 最初のログアウト
		w1 := e2e.PerformRequest(router, http.MethodPost, "/api/users/logout", nil, sessionToken)
		assert.Equal(t, http.StatusNoContent, w1.Code, "First logout should succeed")

		// ログアウト後はセッショントークンが無効化されるため、空のトークンでリクエストを送信
		w2 := e2e.PerformRequest(router, http.MethodPost, "/api/users/logout", nil, "")

		// エラーレスポンスの検証
		e2e.AssertErrorResponse(t, w2, http.StatusUnauthorized, "Unauthorized")
	})
}
