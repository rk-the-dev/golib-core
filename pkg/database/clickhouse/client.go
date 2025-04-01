package clickhouse

import (
	"context"
)

type ClickHouseClient interface {
	Ping(ctx context.Context) error
	Query(ctx context.Context, query string, args ...any) ([]map[string]any, error)
	Exec(ctx context.Context, query string, args ...any) error
	Close() error
}
