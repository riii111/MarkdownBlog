package middleware

import (
	"os"

	"github.com/gin-gonic/contrib/secure"
	"github.com/gin-gonic/gin"
)

// NewSecurityMiddleware セキュリティミドルウェアを初期化
func NewSecurityMiddleware() gin.HandlerFunc {
	return secure.Secure(secure.Options{
		// SSL/HTTPS設定
		SSLRedirect: os.Getenv("SECURITY_SSL_REDIRECT") == "true",
		SSLHost:     os.Getenv("SECURITY_SSL_HOST"),

		// 固定のセキュリティヘッダ
		STSSeconds:           31536000,
		STSIncludeSubdomains: true,
		FrameDeny:            true,
		ContentTypeNosniff:   true,
		BrowserXssFilter:     true,

		// CSPポリシー
		ContentSecurityPolicy: getCSPPolicy(),
	})
}

func getCSPPolicy() string {
	if gin.Mode() == gin.DebugMode {
		return "default-src 'self'; " +
			"script-src 'self' 'unsafe-eval' 'unsafe-inline'; " +
			"style-src 'self' 'unsafe-inline'; " +
			"img-src 'self' data:; " +
			"font-src 'self' data:; " +
			"connect-src 'self' ws: wss:; " +
			"base-uri 'self'; " +
			"form-action 'self'; " +
			"frame-ancestors 'none'; " +
			"object-src 'none'"
	}

	return "default-src 'self'; " +
		"script-src 'self'; " +
		"style-src 'self'; " +
		"img-src 'self' data:; " +
		"font-src 'self'; " +
		"connect-src 'self'; " +
		"base-uri 'self'; " +
		"form-action 'self'; " +
		"frame-ancestors 'none'; " +
		"upgrade-insecure-requests; " +
		"block-all-mixed-content"
}
