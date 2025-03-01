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
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/riii111/markdown-blog-api/internal/handler"
	"github.com/riii111/markdown-blog-api/internal/handler/dto"
	"github.com/riii111/markdown-blog-api/internal/handler/endpoint"
	"github.com/riii111/markdown-blog-api/internal/infrastructure/database"
	"github.com/riii111/markdown-blog-api/internal/infrastructure/migration"
	"github.com/riii111/markdown-blog-api/internal/usecase"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/gorm"
)

// テスト用のデータベース設定
func setupTestDB(t *testing.T) (*gorm.DB, func(), error) {
	ctx := context.Background()

	// PostgreSQLコンテナの設定
	postgresContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15-alpine"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForAll(
				// ログメッセージを確認
				wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
				// ポートが利用可能か確認
				wait.ForListeningPort("5432/tcp"),
			).WithStartupTimeout(10*time.Second),
		),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to start postgres container: %w", err)
	}

	// コンテナのホストとポートを取得
	host, err := postgresContainer.Host(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get container host: %w", err)
	}

	port, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get container port: %w", err)
	}

	// データベース設定
	config := &database.Config{
		Host:     host,
		Port:     port.Port(),
		DBName:   "testdb",
		User:     "testuser",
		Password: "testpass",
	}

	// データベースに接続
	db, err := database.NewDB(config)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// テスト用にAPP_ENV環境変数を設定してマイグレーションを実行
	os.Setenv("APP_ENV", "dev")
	if err := migration.Migrate(db); err != nil {
		return nil, nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	// クリーンアップ関数
	cleanup := func() {
		// テスト終了後にコンテナを停止
		if err := postgresContainer.Terminate(ctx); err != nil {
			log.Printf("Failed to terminate container: %v", err)
		}
	}

	return db, cleanup, nil
}

// テスト用のルーター設定を取得
func getTestRouterConfig() handler.RouterConfig {
	// テスト用の設定を返す
	return handler.TestRouterConfig()
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
	db, dbCleanup, err := setupTestDB(t)
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

	// テスト用のルーターをセットアップ
	router := handler.SetupRouter(userHandler, articleHandler, getTestRouterConfig())

	// クリーンアップ関数
	cleanup := func() {
		// テストデータをクリーンアップ
		cleanupTestData(db)
		// テスト終了後にDBコンテナを停止
		dbCleanup()
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
		// テスト用のセッション名を使用
		if cookie.Name == "test-session" {
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
			Name:  "test-session",
			Value: sessionToken,
		})
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

// テスト用のコンテキストを作成
func CreateTestContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}

// テストデータのクリーンアップ
func cleanupTestData(db *gorm.DB) {
	// トランザクションを使用してデータをクリーンアップ
	db.Transaction(func(tx *gorm.DB) error {
		// 外部キー制約を考慮して削除順序を設定
		tx.Exec("DELETE FROM articles")
		tx.Exec("DELETE FROM sessions")
		tx.Exec("DELETE FROM users")
		return nil
	})
}
