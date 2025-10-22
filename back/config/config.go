package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// DatabaseType represents the configured database driver.
type DatabaseType string

const (
	// DBTypeSQLite uses an embedded sqlite database stored on disk.
	DBTypeSQLite DatabaseType = "sqlite"
	// DBTypeMySQL uses a MySQL-compatible DSN for connectivity.
	DBTypeMySQL DatabaseType = "mysql"
	// DBTypePostgres uses a PostgreSQL DSN for connectivity.
	DBTypePostgres DatabaseType = "postgres"
)

// DatabaseConfig collects the inputs required to create a Gorm connection.
type DatabaseConfig struct {
	Type DatabaseType
	DSN  string
}

// Config centralises configuration used by the application runtime.
type Config struct {
	Port      string
	StaticDir string
	Database  DatabaseConfig
}

// Load reads environment variables and provides sane defaults so the
// application can start without additional setup.
func Load() Config {
	port := firstNonEmpty(os.Getenv("SERVER_PORT"), "8080")

	// default static directory is the dist folder inside the repository.
	cwd, _ := os.Getwd()
	staticDir := firstNonEmpty(os.Getenv("STATIC_DIR"), filepath.Join(cwd, "webserver", "dist"))

	dbType := DatabaseType(firstNonEmpty(os.Getenv("DB_TYPE"), string(DBTypeSQLite)))
	dbDSN := os.Getenv("DB_DSN")

	if dbType == DBTypeSQLite {
		if dbDSN == "" {
			dataDir := filepath.Join(cwd, "data")
			_ = os.MkdirAll(dataDir, 0o755)
			dbDSN = filepath.Join(dataDir, "app.db")
		}
		// normalise sqlite DSN to file path syntax
		dbDSN = fmt.Sprintf("file:%s?_busy_timeout=5000&cache=shared", filepath.ToSlash(dbDSN))
	}

	return Config{
		Port:      port,
		StaticDir: staticDir,
		Database: DatabaseConfig{
			Type: dbType,
			DSN:  dbDSN,
		},
	}
}

func (c Config) Address() string {
	return fmt.Sprintf(":%s", c.Port)
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
