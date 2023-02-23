package storage

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"

	"github.com/goccy/go-json"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/mock"
)

// MockDBExecutor sqlx.DB replacement for testing.
type MockDBExecutor struct {
	mock.Mock

	expectedError error
	expectedRows  *sqlx.Rows
	expectedRow   *sqlx.Row
	expectedDest  interface{}

	affectedRowsResult sql.Result
}

func (m *MockDBExecutor) GetContext(_ context.Context, dest interface{}, _ string, _ ...interface{}) error {
	// hack for copying public fields
	data, err := json.Marshal(m.expectedDest)
	if err != nil {
		return fmt.Errorf("error copying expectedDest value: %w", err)
	}

	if err := json.Unmarshal(data, dest); err != nil { // dest expected to be a pointer
		return fmt.Errorf("error copying expectedDest value to dest: %w", err)
	}

	return m.expectedError
}

func (m *MockDBExecutor) NamedExecContext(context.Context, string, interface{}) (sql.Result, error) {
	return m.affectedRowsResult, m.expectedError
}

// DriverName returns "mock".
func (m *MockDBExecutor) DriverName() string {
	return "mock"
}

// Rebind calls to sqlx.Rebind with UNKNOWN bindType.
func (m *MockDBExecutor) Rebind(s string) string {
	return sqlx.Rebind(sqlx.UNKNOWN, s)
}

// BindNamed calls to sqx.Named.
func (m *MockDBExecutor) BindNamed(s string, i interface{}) (string, []interface{}, error) {
	return sqlx.Named(s, i) //nolint:wrapcheck
}

// QueryContext returns expectedRows and expectedError.
func (m *MockDBExecutor) QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error) {
	var rows *sql.Rows

	if m.expectedRows != nil {
		rows = m.expectedRows.Rows
	}

	return rows, m.expectedError
}

// QueryxContext returns expectedRows and expectedError.
func (m *MockDBExecutor) QueryxContext(context.Context, string, ...interface{}) (*sqlx.Rows, error) {
	return m.expectedRows, m.expectedError
}

// QueryRowxContext returns expectedRow.
func (m *MockDBExecutor) QueryRowxContext(context.Context, string, ...interface{}) *sqlx.Row {
	return m.expectedRow
}

// ExecContext returns affectedRowsResult and expectedError.
func (m *MockDBExecutor) ExecContext(context.Context, string, ...interface{}) (sql.Result, error) {
	return m.affectedRowsResult, m.expectedError
}

func (m *MockDBExecutor) WithError(err error) *MockDBExecutor {
	m.expectedError = err

	return m
}

func (m *MockDBExecutor) WithExpectedResult(result sql.Result) *MockDBExecutor {
	m.affectedRowsResult = result

	return m
}

func (m *MockDBExecutor) WithExpectedRows(rows *sqlx.Rows) *MockDBExecutor {
	m.expectedRows = rows

	return m
}

func (m *MockDBExecutor) WithExpectedRow(row *sqlx.Row) *MockDBExecutor {
	m.expectedRow = row

	return m
}

func (m *MockDBExecutor) WithAffectedRowsResult(affected driver.RowsAffected) {
	m.affectedRowsResult = affected
}
