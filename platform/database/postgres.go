package database

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Noviiich/todo-fiber-pgx/pkg/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgreSQLConnection func for connection to PostgreSQL database.
func PostgreSQLConnection() (*pgxpool.Pool, error) {
	// Define database connection settings.
	maxConn, _ := strconv.Atoi(os.Getenv("DB_MAX_CONNECTIONS"))
	maxIdleConn, _ := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONNECTIONS"))
	maxLifetimeConn, _ := strconv.Atoi(os.Getenv("DB_MAX_LIFETIME_CONNECTIONS"))

	// Build PostgreSQL connection URL.
	postgresConnURL, err := utils.ConnectionURLBuilder("postgres")
	if err != nil {
		return nil, err
	}

	// Parse database URL
	config, err := pgxpool.ParseConfig(postgresConnURL)
	if err != nil {
		return nil, fmt.Errorf("error parsing database URL: %w", err)
	}

	// Set database connection settings:
	// 	- SetMaxOpenConns: the default is 0 (unlimited)
	// 	- SetMaxIdleConns: defaultMaxIdleConns = 2
	// 	- SetConnMaxLifetime: 0, connections are reused forever
	config.MaxConns = int32(maxConn)
	config.MinConns = int32(maxIdleConn)
	config.MaxConnLifetime = time.Duration(maxLifetimeConn) * time.Second

	// Create connection pool
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("error creating connection pool: %w", err)
	}

	return pool, nil
}
