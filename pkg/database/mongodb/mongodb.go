package database

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rk-the-dev/golib-core/pkg/logger"
	"github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDBConfig holds the configuration for MongoDB connection
type MongoDBConfig struct {
	URI      string
	Username string
	Password string
	Database string
	Timeout  time.Duration
}

// MongoDBClient defines the interface for MongoDB operations
type MongoDBClient interface {
	GetCollection(collectionName string) *mongo.Collection
	Close() error
}

// mongoDBClientImpl is the concrete implementation of MongoDBClient
type mongoDBClientImpl struct {
	client   *mongo.Client
	database *mongo.Database
}

var (
	instance MongoDBClient
	once     sync.Once
)

// GetMongoDBClient returns a singleton MongoDB client instance
func GetMongoDBClient(cfg MongoDBConfig) (MongoDBClient, error) {
	var err error
	once.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
		defer cancel()
		clientOptions := options.Client().ApplyURI(cfg.URI)
		if cfg.Username != "" && cfg.Password != "" {
			clientOptions.SetAuth(options.Credential{
				Username: cfg.Username,
				Password: cfg.Password,
			})
		}
		logger.Info("Connecting to MongoDB", logrus.Fields{"uri": cfg.URI, "database": cfg.Database})
		client, connErr := mongo.Connect(ctx, clientOptions)
		if connErr != nil {
			err = fmt.Errorf("failed to connect to MongoDB: %w", connErr)
			logger.Error("MongoDB connection failed", logrus.Fields{"error": connErr})
			return
		}
		db := client.Database(cfg.Database)
		instance = &mongoDBClientImpl{
			client:   client,
			database: db,
		}
		logger.Info("Connected to MongoDB successfully", logrus.Fields{"database": cfg.Database})
	})
	return instance, err
}

// GetCollection returns a MongoDB collection
func (m *mongoDBClientImpl) GetCollection(collectionName string) *mongo.Collection {
	logger.Info("Fetching collection", logrus.Fields{"collection": collectionName})
	return m.database.Collection(collectionName)
}

// Close disconnects the MongoDB client
func (m *mongoDBClientImpl) Close() error {
	logger.Info("Closing MongoDB connection", logrus.Fields{"OPS": "DB Close"})
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := m.client.Disconnect(ctx)
	if err != nil {
		logger.Error("Error closing MongoDB connection", logrus.Fields{"error": err})
	}
	return err
}
