package main

import (
	"log"
	"os"

	_ "github.com/riii111/markdown-blog-api/docs"
	"github.com/riii111/markdown-blog-api/internal/handler"
	"github.com/riii111/markdown-blog-api/internal/infrastructure/database"
	"github.com/riii111/markdown-blog-api/internal/infrastructure/migration"
	"github.com/riii111/markdown-blog-api/internal/usecase"
	"gorm.io/gorm"
)

func initDB() (*gorm.DB, error) {
	config := &database.Config{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		DBName:   os.Getenv("POSTGRES_DB"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
	}

	return database.NewDB(config)
}

// @title           Markdown Blog API
// @version         1.0
// @description     This is Markdown Blog API
// @host           localhost:8088
// @BasePath       /
func main() {
	// バリデーションの初期化
	validationRegistry, err := handler.NewValidationRegistry()
	if err != nil {
		log.Fatalf("Failed to create validation registry: %v", err)
	}
	if err := validationRegistry.RegisterCustomValidations(); err != nil {
		log.Fatalf("Failed to register custom validations: %v", err)
	}

	// データベース接続の初期化
	db, err := initDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// マイグレーションの実行
	if err := migration.Migrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// TODO: インスタンス生成ロジックは一式置き換える（動作確認優先）
	// リポジトリの初期化
	userRepo := database.NewUserRepository(db)
	sessionRepo := database.NewSessionRepository(db)

	// ユースケースの初期化
	userUsecase := usecase.NewUserUsecase(userRepo, sessionRepo)

	// ハンドラーの初期化
	userHandler := handler.NewUserHandler(userUsecase)

	// ルーターのセットアップ
	router := handler.SetupRouter(userHandler)

	// サーバーの起動
	log.Println("Server starting on :8088")
	if err := router.Run(":8088"); err != nil {
		log.Fatal(err)
	}
}
