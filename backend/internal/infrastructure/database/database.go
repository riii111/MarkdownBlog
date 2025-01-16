package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Host     string
	Port     string
	DBName   string
	User     string
	Password string
}

func NewDB(config *Config) (*gorm.DB, error) {
	// 必須パラメータのバリデーション
	if config.Host == "" || config.Port == "" || config.DBName == "" || config.User == "" {
		return nil, fmt.Errorf("missing required database configuration")
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DBName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}
