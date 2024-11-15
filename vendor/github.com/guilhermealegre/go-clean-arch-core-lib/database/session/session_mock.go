package session

import (
	"context"
	"database/sql"

	"github.com/gocraft/dbr/v2"
	"github.com/stretchr/testify/mock"
)

func NewSessionMock() *SessionMock {
	return &SessionMock{}
}

type SessionMock struct {
	mock.Mock
}

func (s *SessionMock) Close() error {
	args := s.Called()
	return args.Error(0)
}

func (s *SessionMock) DB() *sql.DB {
	args := s.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*sql.DB)
}

func (s *SessionMock) Begin() (ITx, error) {
	args := s.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(ITx), args.Error(1)
}

func (s *SessionMock) BeginTx(ctx context.Context, opts *sql.TxOptions) (ITx, error) {
	args := s.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(ITx), args.Error(1)
}

func (s *SessionMock) Select(column ...string) *dbr.SelectBuilder {
	var params []interface{}
	for _, c := range column {
		params = append(params, c)
	}

	args := s.Called(params...)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*dbr.SelectBuilder)
}

func (s *SessionMock) SelectBySql(query string, value ...interface{}) *dbr.SelectBuilder {
	args := s.Called(append([]interface{}{query}, value...)...)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*dbr.SelectBuilder)
}

func (s *SessionMock) InsertInto(table string) *dbr.InsertBuilder {
	args := s.Called(table)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*dbr.InsertBuilder)
}

func (s *SessionMock) InsertBySql(query string, value ...interface{}) *dbr.InsertBuilder {
	args := s.Called(append([]interface{}{query}, value...)...)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*dbr.InsertBuilder)
}

func (s *SessionMock) Update(table string) *dbr.UpdateBuilder {
	args := s.Called(table)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*dbr.UpdateBuilder)
}

func (s *SessionMock) UpdateBySql(query string, value ...interface{}) *dbr.UpdateBuilder {
	args := s.Called(append([]interface{}{query}, value...)...)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*dbr.UpdateBuilder)
}

func (s *SessionMock) DeleteFrom(table string) *dbr.DeleteBuilder {
	args := s.Called(table)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*dbr.DeleteBuilder)
}

func (s *SessionMock) DeleteBySql(query string, value ...interface{}) *dbr.DeleteBuilder {
	args := s.Called(append([]interface{}{query}, value...)...)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*dbr.DeleteBuilder)
}
