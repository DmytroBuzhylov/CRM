package provider

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type DataProvider interface {
	SQLConnector
	BeginTx(ctx context.Context) (pgx.Tx, error)
}
