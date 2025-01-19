package handler

import (
	"net/http"
	"os"
	"time"

	"log"

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

	session, err := h.userUsecase.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "Invalid credentials",
		})
		return
	}

	// 環境変数からセッション期間を取得（時間単位の文字列をパース）
	duration, err := time.ParseDuration(os.Getenv("SESSION_DURATION"))
	if err != nil {
		log.Printf("Failed to parse SESSION_DURATION: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Internal server error",
		})
		return
	}

	// Cookie設定の値をログ出力
	log.Printf("Setting cookie with values: name=%s, token=%s, maxAge=%d, path=%s, domain=%s, secure=%v, httpOnly=%v",
		os.Getenv("SESSION_COOKIE_NAME"),
		session.SessionToken,
		int(duration.Seconds()),
		os.Getenv("COOKIE_PATH"),
		os.Getenv("COOKIE_DOMAIN"),
		os.Getenv("COOKIE_SECURE") == "true",
		os.Getenv("COOKIE_HTTP_ONLY") == "true",
	)

	// セッションCookieを設定
	c.SetCookie(
		os.Getenv("SESSION_COOKIE_NAME"),
		session.SessionToken,
		int(duration.Seconds()), // 秒単位に変換
		os.Getenv("COOKIE_PATH"),
		os.Getenv("COOKIE_DOMAIN"),
		os.Getenv("COOKIE_SECURE") == "true",
		os.Getenv("COOKIE_HTTP_ONLY") == "true",
	)

	// レスポンス
	c.JSON(http.StatusOK, dto.LoginResponse{
		ID:          session.UserID,
		DisplayName: session.User.DisplayName,
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
	sessionToken, err := c.Cookie(os.Getenv("SESSION_COOKIE_NAME"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "No active session",
		})
		return
	}

	if err := h.userUsecase.Logout(c.Request.Context(), sessionToken); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to logout",
		})
		return
	}

	// セッションCookieを削除
	c.SetCookie(
		os.Getenv("SESSION_COOKIE_NAME"),
		"",
		-1,
		os.Getenv("COOKIE_PATH"),
		os.Getenv("COOKIE_DOMAIN"),
		os.Getenv("COOKIE_SECURE") == "true",
		os.Getenv("COOKIE_HTTP_ONLY") == "true",
	)

	c.Status(http.StatusNoContent)
}
