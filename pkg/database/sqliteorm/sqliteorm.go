package sqliteorm

import (
	"fmt"
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SQLiteConfig defines SQLite database configurations
type SQLiteConfig struct {
	DatabaseFile string `env:"SQLITE_DB_FILE" envDefault:"database.db"`
	LogLevel     string `env:"SQLITE_LOG_LEVEL" envDefault:"info"`
}

// SQLiteClient is an implementation of a singleton SQLite database connection
type SQLiteClient struct {
	db *gorm.DB
}

var (
	instance *SQLiteClient
	once     sync.Once
)

// NewSQLiteClient initializes and returns an SQLite database connection
func NewSQLiteClient(cfg *SQLiteConfig) (*SQLiteClient, error) {
	var dbError error
	once.Do(func() {
		// Set GORM log level
		var gormLogLevel logger.LogLevel
		switch cfg.LogLevel {
		case "silent":
			gormLogLevel = logger.Silent
		case "error":
			gormLogLevel = logger.Error
		case "warn":
			gormLogLevel = logger.Warn
		default:
			gormLogLevel = logger.Info
		}
		// Open SQLite connection with GORM default logger
		connection, err := gorm.Open(sqlite.Open(cfg.DatabaseFile), &gorm.Config{
			Logger: logger.Default.LogMode(gormLogLevel),
		})
		if err != nil {
			dbError = fmt.Errorf("❌ Failed to connect to SQLite: %v", err)
			fmt.Println(dbError)
			return
		}
		fmt.Println("✅ Connected to SQLite database:", cfg.DatabaseFile)
		instance = &SQLiteClient{db: connection}
	})
	return instance, dbError
}

// GetDB returns the GORM database instance
func (c *SQLiteClient) GetDB() *gorm.DB {
	return c.db
}
