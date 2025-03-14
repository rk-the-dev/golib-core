package mysqlquery

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// DBClient defines the interface for executing queries
type DBClient interface {
	ExecuteQuery(query string, args ...interface{}) (*sql.Rows, error)
	ExecuteNonQuery(query string, args ...interface{}) (sql.Result, error)
	WithTransaction(txFunc func(*sql.Tx) error) error
	Close() error
}

// MySQLQueryConfig defines MySQL database configurations
type MySQLQueryConfig struct {
	Host            string `env:"DB_HOST" envDefault:"localhost"`
	Port            int    `env:"DB_PORT" envDefault:"3306"`
	User            string `env:"DB_USER" envDefault:"root"`
	Password        string `env:"DB_PASSWORD" envDefault:"password"`
	DatabaseName    string `env:"DB_NAME" envDefault:"mydb"`
	MaxOpenConns    int    `env:"DB_MAX_OPEN_CONNS" envDefault:"100"`
	MaxIdleConns    int    `env:"DB_MAX_IDLE_CONNS" envDefault:"10"`
	ConnMaxLifetime int    `env:"DB_CONN_MAX_LIFETIME" envDefault:"5"`
}

// MySQLClient is an implementation of DBClient
type MySQLClient struct {
	db *sql.DB
}

var (
	instance *MySQLClient
	once     sync.Once
)

// NewMySQLClient initializes and returns a MySQL database connection (raw SQL)
func NewMySQLClient(cfg *MySQLQueryConfig) (DBClient, error) {
	var dbError error
	once.Do(func() {
		// Construct DSN
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DatabaseName)
		// Open MySQL connection
		connection, err := sql.Open("mysql", dsn)
		if err != nil {
			dbError = fmt.Errorf("❌ Failed to connect to MySQL: %v", err)
			fmt.Println(dbError)
			return
		}
		// Configure connection pooling
		connection.SetMaxOpenConns(cfg.MaxOpenConns)
		connection.SetMaxIdleConns(cfg.MaxIdleConns)
		connection.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Minute)
		// Ping the database
		if err := connection.Ping(); err != nil {
			dbError = fmt.Errorf("❌ MySQL connection test failed: %v", err)
			fmt.Println(dbError)
			return
		}
		fmt.Println("✅ Connected to MySQL database:", cfg.DatabaseName)
		instance = &MySQLClient{db: connection}
	})

	return instance, dbError
}

// ExecuteQuery runs a query that returns multiple rows
func (c *MySQLClient) ExecuteQuery(query string, args ...interface{}) (*sql.Rows, error) {
	return c.db.Query(query, args...)
}

// ExecuteNonQuery runs a query that does not return rows (e.g., INSERT, UPDATE, DELETE)
func (c *MySQLClient) ExecuteNonQuery(query string, args ...interface{}) (sql.Result, error) {
	return c.db.Exec(query, args...)
}

// WithTransaction wraps multiple queries inside a transaction
func (c *MySQLClient) WithTransaction(txFunc func(*sql.Tx) error) error {
	tx, err := c.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}
	// Execute function within transaction
	err = txFunc(tx)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("transaction rolled back due to error: %v", err)
	}
	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}
	return nil
}

// Close closes the database connection
func (c *MySQLClient) Close() error {
	return c.db.Close()
}
