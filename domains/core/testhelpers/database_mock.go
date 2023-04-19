package testhelpers

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/mock"
)

// MockDBExecutor sqlx.DB replacement for testing.
type MockDBExecutor struct {
	mock.Mock
}

// GetContext mocks sqlx.DB GetContext method, loading m.expectedDest to dest.
func (m *MockDBExecutor) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	mockArgs := m.Called(ctx, dest, query, args)

	return mockArgs.Error(0)
}

// NamedExecContext returns m.affectedRowsResult, m.expectedError.
func (m *MockDBExecutor) NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
	mockArgs := m.Called(ctx, query, arg)

	return mockArgs.Get(0).(sql.Result), mockArgs.Error(1)
}

// DriverName returns "mock".
func (*MockDBExecutor) DriverName() string {
	return "mock"
}

// Rebind calls to sqlx.Rebind with UNKNOWN bindType.
func (m *MockDBExecutor) Rebind(s string) string {
	mockArgs := m.Called(s)

	return mockArgs.String(0)
}

// BindNamed calls to sqx.Named.
func (m *MockDBExecutor) BindNamed(s string, i interface{}) (string, []interface{}, error) {
	mockArgs := m.Called(s, i)

	return mockArgs.String(0), mockArgs.Get(1).([]interface{}), mockArgs.Error(2)
}

// QueryContext returns expectedRows and expectedError.
func (m *MockDBExecutor) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	mockArgs := m.Called(ctx, query, args)

	return mockArgs.Get(0).(*sql.Rows), mockArgs.Error(1)
}

// QueryxContext returns expectedRows and expectedError.
func (m *MockDBExecutor) QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	mockArgs := m.Called(ctx, query, args)

	return mockArgs.Get(0).(*sqlx.Rows), mockArgs.Error(1)
}

// QueryRowxContext returns expectedRow.
func (m *MockDBExecutor) QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row {
	mockArgs := m.Called(ctx, query, args)

	return mockArgs.Get(0).(*sqlx.Row)
}

// ExecContext returns affectedRowsResult and expectedError.
func (m *MockDBExecutor) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	mockArgs := m.Called(ctx, query, args)

	return mockArgs.Get(0).(sql.Result), mockArgs.Error(1)
}
