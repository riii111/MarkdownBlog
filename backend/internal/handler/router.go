package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/riii111/markdown-blog-api/internal/handler/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(userHandler *UserHandler) *gin.Engine {
	r := gin.Default()

	// CORSミドルウェアを設定
	r.Use(middleware.NewCorsMiddleware())

	// Swaggerのルーティング
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// ヘルスチェック
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// APIグループを削除し、直接ルーティング
	users := r.Group("/api/users")
	{
		users.POST("/register", userHandler.Register)
	}

	return r
}
