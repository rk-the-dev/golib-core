package cassandra

import (
	"context"
)

type Config struct {
	Hosts    []string `mapstructure:"hosts"`
	Keyspace string   `mapstructure:"keyspace"`
	Username string   `mapstructure:"username"`
	Password string   `mapstructure:"password"`
}
type CassandraClient interface {
	Exec(ctx context.Context, query string, args ...any) error
	QueryOne(ctx context.Context, query string, args ...any) (map[string]any, error)
	QueryAll(ctx context.Context, query string, args ...any) ([]map[string]any, error)
	Close()
}
