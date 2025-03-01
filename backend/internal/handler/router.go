package handler

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/riii111/markdown-blog-api/internal/handler/endpoint"
	"github.com/riii111/markdown-blog-api/internal/handler/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type RouterConfig struct {
	SessionName    string
	SessionKey     []byte
	CookiePath     string
	CookieDomain   string
	CookieMaxAge   int
	CookieSecure   bool
	IsTestMode     bool
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

// 開発・本番環境用のデフォルト設定
func DefaultRouterConfig() RouterConfig {
	return RouterConfig{
		SessionName:    os.Getenv("SESSION_NAME"),
		SessionKey:     []byte(os.Getenv("SESSION_SECRET")),
		CookiePath:     os.Getenv("COOKIE_PATH"),
		CookieDomain:   os.Getenv("COOKIE_DOMAIN"),
		CookieMaxAge:   86400 * 30, // 30日
		CookieSecure:   os.Getenv("COOKIE_SECURE") == "true",
		IsTestMode:     false,
		AllowedOrigins: strings.Split(os.Getenv("ALLOWED_ORIGINS"), ","),
		AllowedMethods: strings.Split(os.Getenv("ALLOWED_METHODS"), ","),
		AllowedHeaders: strings.Split(os.Getenv("ALLOWED_HEADERS"), ","),
	}
}

// テスト環境用の設定
func TestRouterConfig() RouterConfig {
	return RouterConfig{
		SessionName:    "test-session",
		SessionKey:     []byte("test-session-secret-key-for-testing-only"),
		CookiePath:     "/",
		CookieDomain:   "localhost",
		CookieMaxAge:   86400,
		CookieSecure:   false,
		IsTestMode:     true,
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}
}

// アプリケーションのルーター
// config パラメーターで本番環境とテスト環境の設定を切り替えることが可
func SetupRouter(
	userHandler *endpoint.UserHandler,
	articleHandler *endpoint.ArticleHandler,
	config ...RouterConfig,
) *gin.Engine {
	// デフォルト設定を使用
	routerConfig := DefaultRouterConfig()
	if len(config) > 0 {
		// 指定された設定で上書き
		routerConfig = config[0]
	}

	// テストモードの場合はテストモードに設定
	if routerConfig.IsTestMode {
		gin.SetMode(gin.TestMode)
	}

	r := gin.Default()

	// セッションキーの確認
	sessionKey := routerConfig.SessionKey
	if len(sessionKey) < 32 && !routerConfig.IsTestMode {
		// キーが32バイト未満の場合は32バイトに拡張（本番環境のみ）
		newKey := make([]byte, 32)
		copy(newKey, sessionKey)
		sessionKey = newKey
	}

	// セッションストアの設定
	store := cookie.NewStore(sessionKey)
	store.Options(sessions.Options{
		Path:     routerConfig.CookiePath,
		Domain:   routerConfig.CookieDomain,
		MaxAge:   routerConfig.CookieMaxAge,
		Secure:   routerConfig.CookieSecure,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	// セッションミドルウェアを設定
	r.Use(sessions.Sessions(routerConfig.SessionName, store))

	// セキュリティミドルウェア
	r.Use(middleware.NewSecurityMiddleware())

	// CORS設定
	if routerConfig.IsTestMode {
		// テスト環境用のシンプルなCORS設定
		r.Use(cors.New(cors.Config{
			AllowOrigins:     routerConfig.AllowedOrigins,
			AllowMethods:     routerConfig.AllowedMethods,
			AllowHeaders:     routerConfig.AllowedHeaders,
			AllowCredentials: true,
			MaxAge:           12 * 3600,
		}))
	} else {
		// 本番環境用のCORS設定
		r.Use(middleware.NewCorsMiddleware())
	}

	// タイムアウトミドルウェア
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
	{
		users := protected.Group("/users")
		{
			users.POST("/logout", userHandler.Logout)
		}

		articles := protected.Group("/articles")
		{
			articles.GET("/me", articleHandler.GetMeArticles)
			articles.POST("", articleHandler.CreateArticle)
			articles.DELETE("/:slug", articleHandler.DeleteArticle)
		}
	}

	// 開発環境でのみSwaggerを有効化（テスト環境では無効）
	if !routerConfig.IsTestMode &&
		strings.ToLower(os.Getenv("APP_ENV")) == "dev" &&
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
