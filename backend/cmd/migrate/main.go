package main

import (
	"flag"
	"log"
	"os"

	"github.com/riii111/markdown-blog-api/internal/infrastructure/database"
	"github.com/riii111/markdown-blog-api/internal/infrastructure/migration"
)

func main() {
	// 環境変数からも設定できるようにflagを使用
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

	if err := migration.Migrate(db); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Migration completed successfully")
}
