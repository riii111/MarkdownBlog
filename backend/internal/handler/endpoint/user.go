package endpoint

import (
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/riii111/markdown-blog-api/internal/handler/dto"
	"github.com/riii111/markdown-blog-api/internal/usecase"
)

// TODO: あとで見直し
type ErrorResponse struct {
	Error string `json:"error" example:"error message"`
}

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

	// レスポンス
	c.JSON(http.StatusCreated, dto.RegisterUserResponse{
		ID:          user.ID,
		DisplayName: user.DisplayName,
		CreatedAt:   user.BaseModel.CreatedAt,
	})
}

// Login godoc
// @Summary      Login user
// @Description  Login with email and password
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body dto.LoginRequest true "Login credentials"
// @Success      200  {object}  dto.LoginResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/users/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request format",
		})
		return
	}

	user, err := h.userUsecase.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "Invalid credentials",
		})
		return
	}

	// セッションの取得と検証
	session := sessions.Default(c)

	// セッションデータを設定
	session.Set("user_id", user.ID.String())
	session.Set("email", user.Email)

	// セッションを保存
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to create session",
		})
		return
	}

	c.JSON(http.StatusOK, dto.LoginResponse{
		ID:          user.ID,
		DisplayName: user.DisplayName,
	})
}

// Logout godoc
// @Summary      Logout user
// @Description  Logout current user and invalidate session
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Success      204
// @Failure      401  {object}  ErrorResponse
// @Router       /api/users/logout [post]
func (h *UserHandler) Logout(c *gin.Context) {
	session := sessions.Default(c)

	// セッションの内容をクリア
	session.Clear()

	// セッションCookieを削除するために有効期限を過去に設定
	session.Options(sessions.Options{
		Path:     os.Getenv("COOKIE_PATH"),
		Domain:   os.Getenv("COOKIE_DOMAIN"),
		MaxAge:   -1, // 負の値を設定することでブラウザがすぐにCookieを削除
		Secure:   os.Getenv("COOKIE_SECURE") == "true",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	// 変更を保存
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to logout",
		})
		return
	}

	c.Status(http.StatusNoContent)
}
