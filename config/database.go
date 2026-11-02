package config

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDatabase initializes the GORM database connection with connection pooling
func InitDatabase(cfg *DatabaseConfig) (*gorm.DB, error) {
	var dialector gorm.Dialector

	switch cfg.Driver {
	case "postgres":
		dialector = postgres.Open(cfg.DSN)
	case "sqlite":
		dialector = sqlite.Open(cfg.DSN)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Driver)
	}

	gormLogger := logger.Default.LogMode(logger.Warn)

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger:                 gormLogger,
		PrepareStmt:            true, // Caches prepared statements for max throughput
		SkipDefaultTransaction: true, // Speeds up writes when single query
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve sql.DB: %w", err)
	}

	// Optimize connection pooling
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

	log.Printf("[DB] Successfully connected to database (%s)", cfg.Driver)
	return db, nil
}
