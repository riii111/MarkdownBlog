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

// Register godoc
// @Summary      Register a new user
// @Description  Register a new user with email, password and display name
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body dto.RegisterUserRequest true "User registration data"
// @Success      201  {object}  dto.RegisterUserResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/users/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterUserRequest

	// リクエストボディのbindとvalidation
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	// ユースケースの呼び出し
	user, err := h.userUsecase.Register(c.Request.Context(), req.Email, req.Password, req.DisplayName)
	if err != nil {
		// すべての登録エラーを400 Bad Requestとして扱う
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Registration failed. Please check your input and try again",
		})
		return
	}

	c.JSON(http.StatusCreated, dto.RegisterUserResponse{
		ID:          user.ID,
		DisplayName: user.DisplayName,
		CreatedAt:   user.BaseModel.CreatedAt,
	})
}

type ErrorResponse struct {
	Error string `json:"error" example:"error message"`
}
