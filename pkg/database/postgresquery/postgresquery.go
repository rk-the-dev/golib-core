package postgresquery

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

// DBClient defines the interface for executing queries
type DBClient interface {
	ExecuteQuery(query string, args ...interface{}) (*sql.Rows, error)
	ExecuteNonQuery(query string, args ...interface{}) (sql.Result, error)
	WithTransaction(txFunc func(*sql.Tx) error) error
	Close() error
}

// PostgresQueryConfig defines PostgreSQL database configurations
type PostgresQueryConfig struct {
	Host            string `env:"PG_HOST" envDefault:"localhost"`
	Port            int    `env:"PG_PORT" envDefault:"5432"`
	User            string `env:"PG_USER" envDefault:"postgres"`
	Password        string `env:"PG_PASSWORD" envDefault:"password"`
	DatabaseName    string `env:"PG_NAME" envDefault:"mydb"`
	SSLMode         string `env:"PG_SSL_MODE" envDefault:"disable"`
	MaxOpenConns    int    `env:"PG_MAX_OPEN_CONNS" envDefault:"100"`
	MaxIdleConns    int    `env:"PG_MAX_IDLE_CONNS" envDefault:"10"`
	ConnMaxLifetime int    `env:"PG_CONN_MAX_LIFETIME" envDefault:"5"`
}

// PostgresClient is an implementation of DBClient
type PostgresClient struct {
	db *sql.DB
}

var (
	instance *PostgresClient
	once     sync.Once
)

// NewPostgresClient initializes and returns a PostgreSQL database connection (raw SQL)
func NewPostgresClient(cfg *PostgresQueryConfig) (DBClient, error) {
	var dbError error
	once.Do(func() {
		// Construct DSN
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DatabaseName, cfg.SSLMode)
		// Open PostgreSQL connection
		connection, err := sql.Open("postgres", dsn)
		if err != nil {
			dbError = fmt.Errorf("❌ Failed to connect to PostgreSQL: %v", err)
			fmt.Println(dbError)
			return
		}
		// Configure connection pooling
		connection.SetMaxOpenConns(cfg.MaxOpenConns)
		connection.SetMaxIdleConns(cfg.MaxIdleConns)
		connection.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Minute)
		// Ping the database
		if err := connection.Ping(); err != nil {
			dbError = fmt.Errorf("❌ PostgreSQL connection test failed: %v", err)
			fmt.Println(dbError)
			return
		}
		fmt.Println("✅ Connected to PostgreSQL database:", cfg.DatabaseName)
		instance = &PostgresClient{db: connection}
	})
	return instance, dbError
}

// ExecuteQuery runs a query that returns multiple rows
func (c *PostgresClient) ExecuteQuery(query string, args ...interface{}) (*sql.Rows, error) {
	return c.db.Query(query, args...)
}

// ExecuteNonQuery runs a query that does not return rows (e.g., INSERT, UPDATE, DELETE)
func (c *PostgresClient) ExecuteNonQuery(query string, args ...interface{}) (sql.Result, error) {
	return c.db.Exec(query, args...)
}

// WithTransaction wraps multiple queries inside a transaction
func (c *PostgresClient) WithTransaction(txFunc func(*sql.Tx) error) error {
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
func (c *PostgresClient) Close() error {
	return c.db.Close()
}
