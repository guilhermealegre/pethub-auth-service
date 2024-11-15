package session

import (
	"context"
	"database/sql"

	"github.com/gocraft/dbr/v2"
)

type ISession interface {
	Close() error
	DB() *sql.DB
	Begin() (ITx, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (ITx, error)
	dbr.SessionRunner
}
