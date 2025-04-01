package cassandra

import (
	"context"
	"fmt"
	"time"

	"github.com/gocql/gocql"
)

type client struct {
	session *gocql.Session
}

func New(hosts []string, keyspace, username, password string) (CassandraClient, error) {
	cluster := gocql.NewCluster(hosts...)
	cluster.Keyspace = keyspace
	cluster.Consistency = gocql.Quorum
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: username,
		Password: password,
	}
	cluster.Timeout = 10 * time.Second

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("cassandra connection failed: %w", err)
	}

	return &client{session: session}, nil
}

func (c *client) Exec(ctx context.Context, query string, args ...any) error {
	return c.session.Query(query, args...).WithContext(ctx).Exec()
}

func (c *client) QueryOne(ctx context.Context, query string, args ...any) (map[string]any, error) {
	iter := c.session.Query(query, args...).WithContext(ctx).Iter()
	defer iter.Close()

	row := map[string]any{}
	if !iter.MapScan(row) {
		return nil, fmt.Errorf("no rows found")
	}
	return row, nil
}

func (c *client) QueryAll(ctx context.Context, query string, args ...any) ([]map[string]any, error) {
	iter := c.session.Query(query, args...).WithContext(ctx).Iter()
	defer iter.Close()

	var results []map[string]any
	row := map[string]any{}
	for iter.MapScan(row) {
		copied := make(map[string]any)
		for k, v := range row {
			copied[k] = v
		}
		results = append(results, copied)
		row = map[string]any{}
	}
	return results, nil
}

func (c *client) Close() {
	c.session.Close()
}
