package session

import (
	"time"

	"github.com/gocraft/dbr/v2"
	"github.com/stretchr/testify/mock"
)

func NewTxMock() *TxMock {
	return &TxMock{}
}

type TxMock struct {
	mock.Mock
}

func (t *TxMock) Tx() *dbr.Tx {
	args := t.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*dbr.Tx)
}

func (t *TxMock) GetTimeout() time.Duration {
	args := t.Called()
	return args.Get(0).(time.Duration)
}

func (t *TxMock) Commit() error {
	args := t.Called()
	return args.Error(0)
}

func (t *TxMock) Rollback() error {
	args := t.Called()
	return args.Error(0)
}

func (t *TxMock) RollbackUnlessCommitted() {
	_ = t.Called()
}

func (t *TxMock) Select(column ...string) *dbr.SelectBuilder {
	var params []interface{}
	for _, c := range column {
		params = append(params, c)
	}

	args := t.Called(params...)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*dbr.SelectBuilder)
}

func (t *TxMock) SelectBySql(query string, value ...interface{}) *dbr.SelectBuilder {
	args := t.Called(append([]interface{}{query}, value...)...)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*dbr.SelectBuilder)
}

func (t *TxMock) InsertInto(table string) *dbr.InsertBuilder {
	args := t.Called(table)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*dbr.InsertBuilder)
}

func (t *TxMock) InsertBySql(query string, value ...interface{}) *dbr.InsertBuilder {
	args := t.Called(append([]interface{}{query}, value...)...)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*dbr.InsertBuilder)
}

func (t *TxMock) Update(table string) *dbr.UpdateBuilder {
	args := t.Called(table)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*dbr.UpdateBuilder)
}

func (t *TxMock) UpdateBySql(query string, value ...interface{}) *dbr.UpdateBuilder {
	args := t.Called(append([]interface{}{query}, value...)...)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*dbr.UpdateBuilder)
}

func (t *TxMock) DeleteFrom(table string) *dbr.DeleteBuilder {
	args := t.Called(table)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*dbr.DeleteBuilder)
}

func (t *TxMock) DeleteBySql(query string, value ...interface{}) *dbr.DeleteBuilder {
	args := t.Called(append([]interface{}{query}, value...)...)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*dbr.DeleteBuilder)
}
