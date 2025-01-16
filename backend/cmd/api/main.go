package main

import (
	"log"
	"net/http"
	"os"

	"github.com/riii111/markdown-blog-api/internal/handler"
	"github.com/riii111/markdown-blog-api/internal/infrastructure/database"
	"github.com/riii111/markdown-blog-api/internal/infrastructure/migration"
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

func main() {
	// ヘルスチェックエンドポイントの登録
	http.HandleFunc("/health", handler.HealthCheck)

	// データベース接続の初期化
	db, err := initDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// マイグレーションの実行
	if err := migration.Migrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// サーバーの起動
	log.Println("Server starting on :8088")
	if err := http.ListenAndServe(":8088", nil); err != nil {
		log.Fatal(err)
	}
}
