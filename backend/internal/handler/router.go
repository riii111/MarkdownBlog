package handler

import "github.com/gin-gonic/gin"

func SetupRouter(userHandler *UserHandler) *gin.Engine {
	r := gin.Default()

	// APIグループ
	api := r.Group("/api")
	{
		api.POST("/register", userHandler.Register)
	}

	return r
}
