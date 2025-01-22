package migration

import (
	"fmt"
	"log"
	"os"

	"github.com/riii111/markdown-blog-api/internal/domain/model"
	"gorm.io/gorm"
)

// モデルの一覧を定義
var models = []interface{}{
	&model.User{},
	&model.Session{},
	&model.Series{},
	&model.Article{},
	&model.Tag{},
	&model.ArticleTag{},
}

func Migrate(db *gorm.DB) error {
	if os.Getenv("APP_ENV") != "dev" {
		log.Println("Skipping auto-migration in non-development environment")
		return nil
	}

	log.Println("Running database migrations...")
	err := db.AutoMigrate(models...)
	if err != nil {
		log.Printf("Migration failed: %v", err)
		return err
	}

	log.Println("Migration completed successfully")
	return nil
}

func ShowStatus(db *gorm.DB) error {
	log.Println("Checking migration status...")

	for _, model := range models {
		// テーブル名の取得方法を修正
		tableName := ""
		if t, ok := model.(interface{ TableName() string }); ok {
			tableName = t.TableName()
		} else {
			stmt := &gorm.Statement{DB: db}
			if err := stmt.Parse(model); err != nil {
				return fmt.Errorf("failed to parse model: %v", err)
			}
			tableName = stmt.Schema.Table
		}

		log.Printf("Checking table: %s", tableName) // デバッグ用

		var count int64
		err := db.Raw(`
			SELECT COUNT(*) 
			FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = ?
		`, tableName).Scan(&count).Error
		if err != nil {
			return fmt.Errorf("failed to check table existence: %v", err)
		}

		exists := count > 0

		var columns []struct {
			ColumnName string `gorm:"column:column_name"`
			DataType   string `gorm:"column:data_type"`
			IsNullable string `gorm:"column:is_nullable"`
		}

		if exists {
			if err := db.Raw(`
				SELECT column_name, data_type, is_nullable 
				FROM information_schema.columns 
				WHERE table_schema = 'public' 
				AND table_name = ?
				ORDER BY ordinal_position
			`, tableName).Scan(&columns).Error; err != nil {
				return fmt.Errorf("failed to get column information: %v", err)
			}

			// 結果を表示
			fmt.Printf("\nTable: %s (Exists: %v)\n", tableName, exists)
			fmt.Println("Columns:")
			for _, col := range columns {
				fmt.Printf("  - %s (%s, Nullable: %s)\n",
					col.ColumnName, col.DataType, col.IsNullable)
			}
		} else {
			fmt.Printf("\nTable: %s (Does not exist)\n", tableName)
		}
	}

	return nil
}

func Rollback(db *gorm.DB) error {
	log.Println("Rolling back last migration...")

	// 安全のため、開発環境でのみ実行可能
	if os.Getenv("APP_ENV") != "dev" {
		return fmt.Errorf("rollback is only allowed in development environment")
	}

	// 逆順でテーブルを削除（依存関係を考慮）
	for i := len(models) - 1; i >= 0; i-- {
		model := models[i]
		tableName := db.Model(model).Statement.Table

		// テーブルが存在するか確認
		var count int64
		err := db.Raw(`
			SELECT COUNT(*) 
			FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = ?
		`, tableName).Scan(&count).Error
		if err != nil {
			return fmt.Errorf("failed to check table existence: %v", err)
		}

		if count > 0 {
			// テーブルを削除
			if err := db.Migrator().DropTable(model); err != nil {
				return fmt.Errorf("failed to drop table %s: %v", tableName, err)
			}
			log.Printf("Dropped table: %s", tableName)
		}
	}

	log.Println("Rollback completed successfully")
	return nil
}
