package session

import (
	"time"

	"github.com/gocraft/dbr/v2"
)

type ITx interface {
	Tx() *dbr.Tx
	GetTimeout() time.Duration
	Commit() error
	Rollback() error
	RollbackUnlessCommitted()
	Select(column ...any) *dbr.SelectBuilder
	SelectBySql(query string, value ...interface{}) *dbr.SelectBuilder
	InsertInto(table string) *dbr.InsertBuilder
	InsertBySql(query string, value ...interface{}) *dbr.InsertBuilder
	Update(table string) *dbr.UpdateBuilder
	UpdateBySql(query string, value ...interface{}) *dbr.UpdateBuilder
	DeleteFrom(table string) *dbr.DeleteBuilder
	DeleteBySql(query string, value ...interface{}) *dbr.DeleteBuilder
}

func newTx(tx *dbr.Tx) *Tx {
	return &Tx{
		conn: tx,
	}
}

type Tx struct {
	conn *dbr.Tx
}

func (t *Tx) Tx() *dbr.Tx {
	return t.conn
}

func (t *Tx) GetTimeout() time.Duration {
	return t.Tx().GetTimeout()
}

func (t *Tx) Commit() error {
	return t.Tx().Commit()
}

func (t *Tx) Rollback() error {
	return t.Tx().Rollback()
}

func (t *Tx) RollbackUnlessCommitted() {
	t.Tx().RollbackUnlessCommitted()
}

func (t *Tx) Select(column ...any) *dbr.SelectBuilder {
	return t.Tx().Select(column...)
}

func (t *Tx) SelectBySql(query string, value ...interface{}) *dbr.SelectBuilder {
	return t.Tx().SelectBySql(query, value...)
}

func (t *Tx) InsertInto(table string) *dbr.InsertBuilder {
	return t.Tx().InsertInto(table)
}

func (t *Tx) InsertBySql(query string, value ...interface{}) *dbr.InsertBuilder {
	return t.Tx().InsertBySql(query, value...)
}

func (t *Tx) Update(table string) *dbr.UpdateBuilder {
	return t.Tx().Update(table)
}

func (t *Tx) UpdateBySql(query string, value ...interface{}) *dbr.UpdateBuilder {
	return t.Tx().UpdateBySql(query, value...)
}

func (t *Tx) DeleteFrom(table string) *dbr.DeleteBuilder {
	return t.Tx().DeleteFrom(table)
}

func (t *Tx) DeleteBySql(query string, value ...interface{}) *dbr.DeleteBuilder {
	return t.Tx().DeleteBySql(query, value...)
}
