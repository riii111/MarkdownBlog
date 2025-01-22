package middleware

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
		session, err := sessionRepo.FindByToken(sessionToken)
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
			// 期限切れセッションを非同期で削除
			go func(sessionID uuid.UUID) {
				if err := sessionRepo.Delete(sessionID); err != nil {
					log.Printf("Failed to delete expired session: %v", err)
				}
			}(session.ID)

			// TODO: 期限切れセッションの一括削除バッチ処理の実装
			// 一定期間経過した期限切れセッションを定期的に削除

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
