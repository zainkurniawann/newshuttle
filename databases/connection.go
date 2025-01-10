package databases

import (
	"context"
	"fmt"
	"shuttle/logger"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var postgresDB *sqlx.DB
var mongoClient *mongo.Client
var once sync.Once

func init() {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func PostgresConnection() (*sqlx.DB, error) {
	once.Do(func() {
		dbURI := "postgres://" + viper.GetString("DB_USER") + ":" + viper.GetString("DB_PASSWORD") + "@" + viper.GetString("DB_HOST") + ":" + viper.GetString("DB_PORT") + "/" + viper.GetString("DB_NAME") + "?sslmode=disable"

		conn, err := sqlx.Connect("postgres", dbURI)
		if err != nil {
			logger.LogFatal(err, "Failed to connect to Postgres", map[string]interface{}{"dbURI": dbURI})
		}

		conn.SetMaxOpenConns(25)
		conn.SetMaxIdleConns(10)
		conn.SetConnMaxLifetime(5 * time.Minute)

		postgresDB = conn
	})

	// Check if the connection is still valid
	if postgresDB == nil {
		return nil, fmt.Errorf("postgres connection is not initialized")
	}

	err := postgresDB.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping Postgres: %w", err)
	}

	return postgresDB, nil
}

func ClosePostgresConnection() error {
	if postgresDB != nil {
		return postgresDB.Close()
	}
	return nil
}

func MongoConnection() (*mongo.Client, error) {
	once.Do(func() {
		clientOptions := options.Client().ApplyURI(viper.GetString("MONGO_URI"))
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var err error
		mongoClient, err = mongo.Connect(ctx, clientOptions)
		if err != nil {
			logger.LogFatal(err, "Failed to connect to MongoDB", map[string]interface{}{"mongoURI": viper.GetString("MONGO_URI")})
		}
	})

	if err := mongoClient.Ping(context.Background(), nil); err != nil {
		return nil, err
	}

	return mongoClient, nil
}