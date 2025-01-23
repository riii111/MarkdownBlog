package handler

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/riii111/markdown-blog-api/internal/handler/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(userHandler *UserHandler, articleHandler *ArticleHandler) *gin.Engine {
	r := gin.Default()

	// セッションの設定（最初に行う）
	sessionKey := []byte(os.Getenv("SESSION_SECRET"))
	if len(sessionKey) < 32 {
		// キーが32バイト未満の場合は32バイトに拡張
		newKey := make([]byte, 32)
		copy(newKey, sessionKey)
		sessionKey = newKey
	}

	store := cookie.NewStore(sessionKey)
	store.Options(sessions.Options{
		Path:     os.Getenv("COOKIE_PATH"),
		Domain:   os.Getenv("COOKIE_DOMAIN"),
		MaxAge:   86400 * 30, // 30日
		Secure:   os.Getenv("COOKIE_SECURE") == "true",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	// セッションミドルウェアを最初に設定
	r.Use(sessions.Sessions(os.Getenv("SESSION_NAME"), store))

	// その他のミドルウェア
	r.Use(middleware.NewSecurityMiddleware())
	r.Use(middleware.NewCorsMiddleware())
	r.Use(middleware.TimeoutMiddleware())

	// 認証不要のエンドポイント
	public := r.Group("/api")
	{
		users := public.Group("/users")
		{
			users.POST("/register", userHandler.Register)
			users.POST("/login", userHandler.Login)
		}

		articles := public.Group("/articles")
		{
			articles.GET("", articleHandler.GetArticles)
			articles.GET("/:slug", articleHandler.GetArticleBySlug)
		}
	}

	// 認証が必要なエンドポイント
	protected := r.Group("/api")
	protected.Use(middleware.CSRF())
	protected.Use(middleware.AuthMiddleware())
	protected.Use(middleware.CSRF())
	{
		users := protected.Group("/users")
		{
			users.POST("/logout", userHandler.Logout)
		}

		articles := protected.Group("/articles")
		{
			articles.POST("", articleHandler.CreateArticle)
			articles.DELETE("/:slug", articleHandler.DeleteArticle)
		}
	}

	// 開発環境でのみSwaggerを有効化
	if strings.ToLower(os.Getenv("APP_ENV")) == "dev" &&
		strings.ToLower(os.Getenv("ENABLE_SWAGGER")) == "true" {
		// SwaggerルートにだけBasic認証を適用
		swaggerAuth := r.Group("/swagger", gin.BasicAuth(gin.Accounts{
			os.Getenv("SWAGGER_USER"): os.Getenv("SWAGGER_PASSWORD"),
		}))

		swaggerAuth.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// ヘルスチェック
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	return r
}
