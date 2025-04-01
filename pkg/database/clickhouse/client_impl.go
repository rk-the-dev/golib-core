package clickhouse

import (
	"context"
	"fmt"

	ch "github.com/ClickHouse/clickhouse-go/v2"
)

type client struct {
	conn ch.Conn
}

// New creates a new ClickHouse client using a DSN-like address (e.g., "localhost:9000")
func New(addr, username, password, database string) (ClickHouseClient, error) {
	conn, err := ch.Open(&ch.Options{
		Addr: []string{addr},
		Auth: ch.Auth{
			Database: database,
			Username: username,
			Password: password,
		},
		Settings: ch.Settings{
			"max_execution_time": 60,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("clickhouse connection failed: %w", err)
	}

	return &client{conn: conn}, nil
}

func (c *client) Ping(ctx context.Context) error {
	return c.conn.Ping(ctx)
}

func (c *client) Exec(ctx context.Context, query string, args ...any) error {
	return c.conn.Exec(ctx, query, args...)
}

func (c *client) Query(ctx context.Context, query string, args ...any) ([]map[string]any, error) {
	rows, err := c.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns := rows.Columns()
	colCount := len(columns)

	var results []map[string]any
	for rows.Next() {
		values := make([]any, colCount)
		if err := rows.Scan(values...); err != nil {
			return nil, err
		}
		row := make(map[string]any)
		for i, col := range columns {
			row[col] = values[i]
		}
		results = append(results, row)
	}

	return results, nil
}

func (c *client) Close() error {
	return c.conn.Close()
}
