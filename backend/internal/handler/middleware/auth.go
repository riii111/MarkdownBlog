package middleware

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/riii111/markdown-blog-api/internal/domain/model"
)

func AuthMiddleware(sessionRepo model.SessionRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// セッショントークンをCookieから取得
		sessionToken, err := c.Cookie(os.Getenv("SESSION_COOKIE_NAME"))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No session token found"})
			c.Abort()
			return
		}

		// セッションの検証
		session, err := sessionRepo.FindByToken(c.Request.Context(), sessionToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session"})
			c.Abort()
			return
		}
		if session == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Session not found"})
			c.Abort()
			return
		}

		// 有効期限チェック
		if time.Now().After(session.ExpiresAt) {
			// 期限切れセッションの削除
			if err := sessionRepo.Delete(session.ID); err != nil {
				log.Printf("Failed to delete expired session: %v", err)
				// 削除に失敗しても認証エラーは返す
			}

			// Cookieの削除
			c.SetCookie(
				os.Getenv("SESSION_COOKIE_NAME"),
				"",
				-1,
				os.Getenv("COOKIE_PATH"),
				os.Getenv("COOKIE_DOMAIN"),
				os.Getenv("COOKIE_SECURE") == "true",
				os.Getenv("COOKIE_HTTP_ONLY") == "true",
			)

			c.JSON(http.StatusUnauthorized, gin.H{"error": "Session expired"})
			c.Abort()
			return
		}

		// 有効なセッションの場合はユーザーIDを設定
		c.Set("userID", session.UserID)
		c.Next()
	}
}
