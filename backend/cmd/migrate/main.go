package main

import (
	"flag"
	"log"
	"os"

	"github.com/riii111/markdown-blog-api/internal/infrastructure/database"
	"github.com/riii111/markdown-blog-api/internal/infrastructure/migration"
)

func main() {
	// フラグの定義
	status := flag.Bool("status", false, "Show migration status")
	rollback := flag.Bool("rollback", false, "Rollback last migration")

	dbHost := flag.String("host", os.Getenv("POSTGRES_HOST"), "Database host")
	dbPort := flag.String("port", os.Getenv("POSTGRES_PORT"), "Database port")
	dbName := flag.String("dbname", os.Getenv("POSTGRES_DB"), "Database name")
	dbUser := flag.String("user", os.Getenv("POSTGRES_USER"), "Database user")
	dbPass := flag.String("password", os.Getenv("POSTGRES_PASSWORD"), "Database password")
	flag.Parse()

	config := &database.Config{
		Host:     *dbHost,
		Port:     *dbPort,
		DBName:   *dbName,
		User:     *dbUser,
		Password: *dbPass,
	}

	db, err := database.NewDB(config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 接続確認のログを追加
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get underlying sql.DB: %v", err)
	}

	err = sqlDB.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Successfully connected to database")

	// フラグに応じた処理の分岐
	if *status {
		if err := migration.ShowStatus(db); err != nil {
			log.Fatalf("Failed to show migration status: %v", err)
		}
		return
	}

	if *rollback {
		if err := migration.Rollback(db); err != nil {
			log.Fatalf("Rollback failed: %v", err)
		}
		return
	}

	if err := migration.Migrate(db); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Migration completed successfully")
}
