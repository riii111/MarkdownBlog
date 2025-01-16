package migration

import (
	"github.com/riii111/markdown-blog-api/internal/domain/model"
	"gorm.io/gorm"
	"log"
)

func Migrate(db *gorm.DB) error {
	log.Println("Running database migrations...")

	// マイグレーションの実行
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
