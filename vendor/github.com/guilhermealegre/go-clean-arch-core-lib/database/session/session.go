package session

import (
	"context"
	"database/sql"

	"github.com/gocraft/dbr/v2"
)

type Session struct {
	conn *dbr.Session
}

func NewSession(conn *dbr.Connection, event dbr.EventReceiver) *Session {
	return &Session{
		conn: &dbr.Session{
			Connection:    conn,
			EventReceiver: event,
		},
	}
}

func (s *Session) Close() error {
	return s.conn.Close()
}

func (s *Session) DB() *sql.DB {
	return s.conn.DB
}

func (s *Session) Begin() (ITx, error) {
	tx, err := s.conn.Begin()
	if err != nil {
		return nil, err
	}

	return newTx(tx), nil
}

func (s *Session) BeginTx(ctx context.Context, opts *sql.TxOptions) (ITx, error) {
	tx, err := s.conn.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}

	return newTx(tx), nil
}

func (s *Session) Select(column ...any) *dbr.SelectBuilder {
	return s.conn.Select(column...)
}

func (s *Session) SelectBySql(query string, value ...interface{}) *dbr.SelectBuilder {
	return s.conn.SelectBySql(query, value...)
}

func (s *Session) InsertInto(table string) *dbr.InsertBuilder {
	return s.conn.InsertInto(table)
}

func (s *Session) InsertBySql(query string, value ...interface{}) *dbr.InsertBuilder {
	return s.conn.InsertBySql(query, value...)
}

func (s *Session) Update(table string) *dbr.UpdateBuilder {
	return s.conn.Update(table)
}

func (s *Session) UpdateBySql(query string, value ...interface{}) *dbr.UpdateBuilder {
	return s.conn.UpdateBySql(query, value...)
}

func (s *Session) DeleteFrom(table string) *dbr.DeleteBuilder {
	return s.conn.DeleteFrom(table)
}

func (s *Session) DeleteBySql(query string, value ...interface{}) *dbr.DeleteBuilder {
	return s.conn.DeleteBySql(query, value...)
}
