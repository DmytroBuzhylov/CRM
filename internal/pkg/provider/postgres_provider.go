package provider

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresProvider struct {
	pool *pgxpool.Pool
}

func NewPostgresProvider(pool *pgxpool.Pool) *PostgresProvider {
	return &PostgresProvider{pool: pool}
}

func (p *PostgresProvider) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return p.pool.QueryRow(ctx, sql, args...)
}

func (p *PostgresProvider) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return p.pool.Query(ctx, sql, args...)
}

func (p *PostgresProvider) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	return p.pool.Exec(ctx, sql, args...)
}

func (p *PostgresProvider) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return p.pool.BeginTx(ctx, pgx.TxOptions{})
}
