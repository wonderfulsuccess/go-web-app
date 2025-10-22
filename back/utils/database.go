package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/wonderfulsuccess/go-web-app/back/config"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDatabase opens a database connection using the configured driver.
func InitDatabase(cfg config.DatabaseConfig) (*gorm.DB, error) {
	var (
		db  *gorm.DB
		err error
	)

	gormCfg := &gorm.Config{Logger: logger.Default.LogMode(logger.Info)}

	switch cfg.Type {
	case config.DBTypeSQLite:
		if err := ensureSQLiteFile(cfg.DSN); err != nil {
			return nil, err
		}
		db, err = gorm.Open(sqlite.Open(cfg.DSN), gormCfg)
	case config.DBTypeMySQL:
		dsn := cfg.DSN
		if dsn == "" {
			dsn = "root:password@tcp(127.0.0.1:3306)/app?charset=utf8mb4&parseTime=True&loc=Local"
		}
		db, err = gorm.Open(mysql.Open(dsn), gormCfg)
	case config.DBTypePostgres:
		dsn := cfg.DSN
		if dsn == "" {
			dsn = "host=127.0.0.1 user=postgres password=postgres dbname=app port=5432 sslmode=disable TimeZone=Asia/Shanghai"
		}
		db, err = gorm.Open(postgres.Open(dsn), gormCfg)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.Type)
	}

	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	return db, nil
}

func ensureSQLiteFile(dsn string) error {
	if dsn == "" {
		return fmt.Errorf("sqlite DSN cannot be empty")
	}

	// strip the file: prefix if present
	trimmed := dsn
	if len(trimmed) > 5 && trimmed[:5] == "file:" {
		trimmed = trimmed[5:]
	}
	if idx := indexRune(trimmed, '?'); idx >= 0 {
		trimmed = trimmed[:idx]
	}
	if trimmed == "" {
		return fmt.Errorf("sqlite DSN did not contain a file path")
	}

	path := filepath.FromSlash(trimmed)
	dir := filepath.Dir(path)
	return os.MkdirAll(dir, 0o755)
}

func indexRune(s string, r rune) int {
	for i, ch := range s {
		if ch == r {
			return i
		}
	}
	return -1
}
