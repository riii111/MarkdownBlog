package migration

import (
	"log"
	"os"

	"github.com/riii111/markdown-blog-api/internal/domain/model"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	// 開発環境でのみ自動マイグレーションを実行
	if os.Getenv("APP_ENV") != "development" {
		log.Println("Skipping auto-migration in non-development environment")
		return nil
	}

	log.Println("Running database migrations...")

	err := db.AutoMigrate(
		&model.User{},
		&model.Session{},
		&model.Series{},
		&model.Post{},
		&model.Tag{},
		&model.PostTag{},
	)

	if err != nil {
		log.Printf("Migration failed: %v", err)
		return err
	}

	log.Println("Migration completed successfully")
	return nil
}
