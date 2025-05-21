package db

import (
	"Test/config"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"time"
)

func InitPostgresSQLDB(cfg config.DatabaseConfig) (*pgxpool.Pool, error) {
	connStr := cfg.DSN
	dbConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DSN: %w", err)
	}
	dbConfig.MaxConns = int32(cfg.MaxOpenConns)
	dbConfig.MinConns = int32(cfg.MaxIdleConns)
	dbConfig.MaxConnLifetime = cfg.ConnMaxLifetime
	dbConfig.MaxConnIdleTime = cfg.ConnMaxIdletime
	dbConfig.HealthCheckPeriod = cfg.HealthCheckInterval

	pool, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = pool.Ping(ctx)

	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	log.Info().Msg("Successfully connected to PostgreSQL!")
	return pool, nil
}
