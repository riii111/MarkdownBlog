package handler

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/riii111/markdown-blog-api/internal/handler/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(userHandler *UserHandler, postHandler *PostHandler) *gin.Engine {
	r := gin.Default()

	// CORSミドルウェアを設定
	r.Use(middleware.NewCorsMiddleware())

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

	users := r.Group("/api/users")
	{
		users.POST("/register", userHandler.Register)
		users.POST("/login", userHandler.Login)
		users.POST("/logout", userHandler.Logout)
	}

	posts := r.Group("/api/posts")
	{
		posts.POST("", postHandler.CreatePost)
		posts.DELETE("/:slug", postHandler.DeletePost)
	}

	return r
}
