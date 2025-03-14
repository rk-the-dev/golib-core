package mysqlorm

import (
	"fmt"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// MySqlORMConfig defines MySQL database configurations
type MySqlORMConfig struct {
	Host            string `env:"DB_HOST" envDefault:"localhost"`
	Port            int    `env:"DB_PORT" envDefault:"3306"`
	User            string `env:"DB_USER" envDefault:"root"`
	Password        string `env:"DB_PASSWORD" envDefault:"password"`
	DatabaseName    string `env:"DB_NAME" envDefault:"mydb"`
	Charset         string `env:"DB_CHARSET" envDefault:"utf8mb4"`
	ParseTime       bool   `env:"DB_PARSE_TIME" envDefault:"true"`
	Loc             string `env:"DB_LOC" envDefault:"Local"`
	MaxOpenConns    int    `env:"DB_MAX_OPEN_CONNS" envDefault:"100"`
	MaxIdleConns    int    `env:"DB_MAX_IDLE_CONNS" envDefault:"10"`
	ConnMaxLifetime int    `env:"DB_CONN_MAX_LIFETIME" envDefault:"5"`
	LogLevel        string `env:"DB_LOG_LEVEL" envDefault:"info"`
}

// Singleton instance
var (
	db   *gorm.DB
	once sync.Once
)

// GetDB initializes and returns a MySQL database connection
func GetDB(cfg *MySqlORMConfig) (*gorm.DB, error) {
	var dbError error
	once.Do(func() {
		// Construct DSN
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DatabaseName, cfg.Charset, cfg.ParseTime, cfg.Loc)
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
		// Open MySQL connection with GORM default logger
		connection, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(gormLogLevel),
		})
		if err != nil {
			dbError = fmt.Errorf("❌ Failed to connect to MySQL: %v", err)
			fmt.Println(dbError)
			return
		}
		db = connection // Assigning the new database connection
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
		fmt.Println("✅ Connected to MySQL database:", cfg.DatabaseName)
	})
	return db, dbError
}
