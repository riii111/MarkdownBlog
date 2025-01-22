package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/riii111/markdown-blog-api/internal/pkg/config"
)

const (
	csrfTokenLength = 32
	csrfTokenKey    = "X-CSRF-Token"
)

func generateCSRFToken() string {
	b := make([]byte, csrfTokenLength)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func CSRF() gin.HandlerFunc {
	return func(c *gin.Context) {
		// セッショントークンの取得
		sessionToken, err := c.Cookie(os.Getenv("SESSION_COOKIE_NAME"))
		if err != nil || sessionToken == "" {
			c.Next()
			return
		}

		if c.Request.Method == "GET" {
			token := generateCSRFToken()

			// CSRFトークンをクッキーとして設定
			c.SetCookie(
				csrfTokenKey,
				token,
				config.GetSessionMaxAge(),
				os.Getenv("COOKIE_PATH"),
				os.Getenv("COOKIE_DOMAIN"),
				os.Getenv("COOKIE_SECURE") == "true",
				true,
			)

			// レスポンスヘッダーにも設定（フロントエンドはこちらを使用する）
			c.Header(csrfTokenKey, token)
		} else {
			// POST, PUT, DELETE等の場合、トークンを検証
			cookieToken, _ := c.Cookie(csrfTokenKey)
			headerToken := c.GetHeader(csrfTokenKey)

			if cookieToken == "" || headerToken == "" || cookieToken != headerToken {
				c.JSON(http.StatusForbidden, gin.H{
					"error": "CSRF token validation failed",
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
