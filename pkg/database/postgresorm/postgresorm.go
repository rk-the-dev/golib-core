package postgresorm

import (
	"fmt"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// PostgresORMConfig defines PostgreSQL database configurations
type PostgresORMConfig struct {
	Host            string `env:"PG_HOST" envDefault:"localhost"`
	Port            int    `env:"PG_PORT" envDefault:"5432"`
	User            string `env:"PG_USER" envDefault:"postgres"`
	Password        string `env:"PG_PASSWORD" envDefault:"password"`
	DatabaseName    string `env:"PG_NAME" envDefault:"mydb"`
	SSLMode         string `env:"PG_SSL_MODE" envDefault:"disable"`
	MaxOpenConns    int    `env:"PG_MAX_OPEN_CONNS" envDefault:"100"`
	MaxIdleConns    int    `env:"PG_MAX_IDLE_CONNS" envDefault:"10"`
	ConnMaxLifetime int    `env:"PG_CONN_MAX_LIFETIME" envDefault:"5"`
	LogLevel        string `env:"PG_LOG_LEVEL" envDefault:"info"`
}

// Singleton instance
var (
	db   *gorm.DB
	once sync.Once
)

// GetDB initializes and returns a PostgreSQL database connection
func GetDB(cfg *PostgresORMConfig) (*gorm.DB, error) {
	var dbError error

	once.Do(func() {
		// Construct DSN
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DatabaseName, cfg.SSLMode)

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

		// Open PostgreSQL connection with GORM default logger
		connection, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(gormLogLevel),
		})
		if err != nil {
			dbError = fmt.Errorf("❌ Failed to connect to PostgreSQL: %v", err)
			fmt.Println(dbError)
			return
		}
		db = connection // Assign the new database connection

		// Configure connection pooling
		sqlDB, err := db.DB()
		if err != nil {
			dbError = fmt.Errorf("❌ Failed to retrieve SQL DB instance: %v", err)
			fmt.Println(dbError)
			db = nil
			return
		}

		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
		sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Minute)

		fmt.Println("✅ Connected to PostgreSQL database:", cfg.DatabaseName)
	})

	return db, dbError
}
