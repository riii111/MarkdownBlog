package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const defaultTimeout = 10 * time.Second

// DBへのリクエストタイムアウトのミドルウェア
func TimeoutMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), defaultTimeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)

		c.Next()

		// コンテキストがタイムアウトした場合のハンドリング
		if ctx.Err() == context.DeadlineExceeded {
			c.AbortWithStatusJSON(http.StatusGatewayTimeout, gin.H{
				"error": "Request timeout",
			})
			return
		}
	}
}
