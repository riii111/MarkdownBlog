package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/riii111/markdown-blog-api/internal/handler/dto"
	"github.com/riii111/markdown-blog-api/internal/usecase"
)

type UserHandler struct {
	userUsecase *usecase.UserUsecase
}
func NewUserHandler(userUsecase *usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
	}
}

// User登録ハンドラ
func (h *UserHandler) Register(c *.gin.Context) {
	var req dto.RegisterUserRequest

	// リクエストボディのbindとvalidation
	if err := c.ShouldBindJSON(&req) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// ユースケースの呼び出し
	user, err := h.userUsecase.Register(c.Request.Context(), req.Email, req.Password, req.DisplayName)
	if err != nil {
		// TODO: あとでエラーハンドリング見直す
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to register user",
		})
		return
	}

	// 成功時
	c.JSON(http.StatusCreated, dto.RegisterUserResponse{
		ID:	user.ID.string(),
		DisplayName:  user.DisplayName,
		CreatedAt:	user.CreatedAt,
	})
}