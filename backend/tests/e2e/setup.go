package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/riii111/markdown-blog-api/internal/handler"
	"github.com/riii111/markdown-blog-api/internal/handler/dto"
	"github.com/riii111/markdown-blog-api/internal/handler/endpoint"
	"github.com/riii111/markdown-blog-api/internal/infrastructure/database"
	"github.com/riii111/markdown-blog-api/internal/infrastructure/migration"
	"github.com/riii111/markdown-blog-api/internal/usecase"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// テスト用のデータベース設定
func setupTestDB() (*gorm.DB, error) {
	// テスト用のデータベース設定
	config := &database.Config{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		DBName:   os.Getenv("TEST_DB_NAME"), // テスト用DBを使用
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
	}

	// 環境変数が設定されていない場合はデフォルト値を使用
	if config.Host == "" {
		config.Host = "localhost"
	}
	if config.Port == "" {
		config.Port = "5432"
	}
	if config.DBName == "" {
		config.DBName = "markdown_blog_test"
	}
	if config.User == "" {
		config.User = "postgres"
	}

	db, err := database.NewDB(config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to test database: %w", err)
	}

	// マイグレーションを実行
	if err := migration.Migrate(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return db, nil
}

// テスト環境のセットアップ
func SetupTestEnvironment(t *testing.T) (*gin.Engine, func()) {
	// テストモードに設定
	gin.SetMode(gin.TestMode)

	// バリデーションの初期化
	validationRegistry, err := handler.NewValidationRegistry()
	require.NoError(t, err, "Failed to create validation registry")
	err = validationRegistry.RegisterCustomValidations()
	require.NoError(t, err, "Failed to register custom validations")

	// テスト用DBの設定
	db, err := setupTestDB()
	require.NoError(t, err, "Failed to setup test database")

	// リポジトリの初期化
	userRepo := database.NewUserRepository(db)
	sessionRepo := database.NewSessionRepository(db)
	articleRepo := database.NewArticleRepository(db)

	// ユースケースの初期化
	userUsecase := usecase.NewUserUsecase(userRepo, sessionRepo)
	articleUsecase := usecase.NewArticleUsecase(articleRepo)

	// ハンドラーの初期化
	userHandler := endpoint.NewUserHandler(userUsecase)
	articleHandler := endpoint.NewArticleHandler(articleUsecase)

	// ルーターのセットアップ
	router := handler.SetupRouter(userHandler, articleHandler)

	// クリーンアップ関数
	cleanup := func() {
		// テスト終了後にテーブルをクリーンアップ
		sqlDB, err := db.DB()
		if err != nil {
			log.Printf("Failed to get DB instance: %v", err)
			return
		}
		sqlDB.Close()
	}

	return router, cleanup
}

// テスト用ユーザーの作成
func CreateTestUser(t *testing.T, router *gin.Engine) (dto.RegisterUserResponse, string) {
	// テスト用ユーザーデータ
	testUser := dto.RegisterUserRequest{
		Email:       fmt.Sprintf("test-%s@example.com", uuid.New().String()),
		Password:    "password123",
		DisplayName: "Test User",
	}

	// ユーザー登録リクエスト
	body, err := json.Marshal(testUser)
	require.NoError(t, err, "Failed to marshal test user")

	req := httptest.NewRequest(http.MethodPost, "/api/users/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code, "Expected status code 201")

	// レスポンスをパース
	var response dto.RegisterUserResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err, "Failed to unmarshal response")

	// ログインしてセッショントークンを取得
	loginReq := dto.LoginRequest{
		Email:    testUser.Email,
		Password: testUser.Password,
	}

	loginBody, err := json.Marshal(loginReq)
	require.NoError(t, err, "Failed to marshal login request")

	loginReqObj := httptest.NewRequest(http.MethodPost, "/api/users/login", bytes.NewBuffer(loginBody))
	loginReqObj.Header.Set("Content-Type", "application/json")

	loginW := httptest.NewRecorder()
	router.ServeHTTP(loginW, loginReqObj)

	require.Equal(t, http.StatusOK, loginW.Code, "Expected status code 200 for login")

	// セッションCookieを取得
	cookies := loginW.Result().Cookies()
	var sessionToken string
	for _, cookie := range cookies {
		if cookie.Name == os.Getenv("SESSION_NAME") {
			sessionToken = cookie.Value
			break
		}
	}

	require.NotEmpty(t, sessionToken, "Session token should not be empty")

	return response, sessionToken
}

// HTTPリクエストを実行するヘルパー関数
func PerformRequest(router *gin.Engine, method, path string, body interface{}, sessionToken string) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(jsonBody)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req := httptest.NewRequest(method, path, reqBody)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// セッショントークンがある場合はCookieを設定
	if sessionToken != "" {
		req.AddCookie(&http.Cookie{
			Name:  os.Getenv("SESSION_NAME"),
			Value: sessionToken,
		})
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

// テスト用のコンテキストを作成
func CreateTestContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5000)
}
