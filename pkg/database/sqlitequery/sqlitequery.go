package sqlitequery

import (
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

// DBClient defines the interface for executing queries
type DBClient interface {
	ExecuteQuery(query string, args ...interface{}) (*sql.Rows, error)
	ExecuteNonQuery(query string, args ...interface{}) (sql.Result, error)
	WithTransaction(txFunc func(*sql.Tx) error) error
	Close() error
}

// SQLiteConfig defines SQLite database configurations
type SQLiteConfig struct {
	DatabaseFile string `env:"SQLITE_DB_FILE" envDefault:"database.db"`
}

// SQLiteClient is an implementation of DBClient
type SQLiteClient struct {
	db *sql.DB
}

var (
	instance *SQLiteClient
	once     sync.Once
)

// NewSQLiteClient initializes and returns an SQLite database connection
func NewSQLiteClient(cfg *SQLiteConfig) (DBClient, error) {
	var dbError error
	once.Do(func() {
		// Open SQLite connection
		connection, err := sql.Open("sqlite3", cfg.DatabaseFile)
		if err != nil {
			dbError = fmt.Errorf("❌ Failed to connect to SQLite: %v", err)
			fmt.Println(dbError)
		}
		// Ping the database
		if err := connection.Ping(); err != nil {
			dbError = fmt.Errorf("❌ SQLite connection test failed: %v", err)
			fmt.Println(dbError)
		}
		fmt.Println("✅ Connected to SQLite database:", cfg.DatabaseFile)
		instance = &SQLiteClient{db: connection}
	})
	return instance, dbError
}

// ExecuteQuery runs a query that returns multiple rows
func (c *SQLiteClient) ExecuteQuery(query string, args ...interface{}) (*sql.Rows, error) {
	return c.db.Query(query, args...)
}

// ExecuteNonQuery runs a query that does not return rows (e.g., INSERT, UPDATE, DELETE)
func (c *SQLiteClient) ExecuteNonQuery(query string, args ...interface{}) (sql.Result, error) {
	return c.db.Exec(query, args...)
}

// WithTransaction wraps multiple queries inside a transaction
func (c *SQLiteClient) WithTransaction(txFunc func(*sql.Tx) error) error {
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
func (c *SQLiteClient) Close() error {
	return c.db.Close()
}
