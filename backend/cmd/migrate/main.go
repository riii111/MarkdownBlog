package main

import (
	"flag"
	"log"

	"github.com/riii111/markdown-blog-api/internal/infrastructure/database"
	"github.com/riii111/markdown-blog-api/internal/infrastructure/migration"
)

func main() {
	// 環境変数からも設定できるようにflagを使用
	dbHost := flag.String("host", "localhost", "Database host")
	dbPort := flag.String("port", "5432", "Database port")
	dbName := flag.String("dbname", "blog_db", "Database name")
	dbUser := flag.String("user", "postgres", "Database user")
	dbPass := flag.String("password", "", "Database password")
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
