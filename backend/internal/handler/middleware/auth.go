package middleware

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/riii111/markdown-blog-api/internal/domain/model"
)

func AuthMiddleware(sessionRepo model.SessionRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// セッショントークンをCookieから取得
		sessionToken, err := c.Cookie(os.Getenv("SESSION_COOKIE_NAME"))
		if err != nil {
			log.Println("session token not found in cookie")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// セッションの検証
		session, err := sessionRepo.FindByToken(c.Request.Context(), sessionToken)
		if err != nil || session == nil {
			log.Println("session not found in DB")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session"})
			c.Abort()
			return
		}

		// セッションの有効期限チェック
		if session.IsExpired() {
			log.Println("session expired")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Session expired"})
			c.Abort()
			return
		}

		// ユーザーIDをコンテキストに設定
		c.Set("userID", session.UserID)
		c.Next()
	}
}
