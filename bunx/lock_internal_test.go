package bunx

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

// fakeResult controls how the fake driver below responds to the advisory-lock
// query, so that tryTakeAdvisoryLock's error branches (query/no-rows/scan) can
// be driven without a real database. A nil row means no rows are returned.
type fakeResult struct {
	queryErr error
	row      []driver.Value
}

// beginFakeTx wires up a bun.Tx backed by a minimal in-memory database/sql/driver
// fake configured by result.
func beginFakeTx(t *testing.T, result fakeResult) bun.Tx {
	t.Helper()

	sqlDB := sql.OpenDB(&fakeConnector{result: result})
	db := bun.NewDB(sqlDB, pgdialect.New())
	t.Cleanup(func() { _ = db.Close() })

	tx, err := db.BeginTx(context.Background(), nil)
	require.NoError(t, err)
	t.Cleanup(func() { _ = tx.Rollback() })

	return tx
}

type fakeConnector struct {
	result fakeResult
}

func (c *fakeConnector) Connect(context.Context) (driver.Conn, error) {
	return &fakeConn{result: c.result}, nil
}

func (c *fakeConnector) Driver() driver.Driver {
	return &fakeDriver{}
}

// fakeDriver only exists to satisfy driver.Connector.Driver(); sql.OpenDB uses
// the connector directly and never calls Open.
type fakeDriver struct{}

func (*fakeDriver) Open(string) (driver.Conn, error) {
	return nil, errors.New("fakeDriver: Open is not supported, use the connector")
}

type fakeConn struct {
	result fakeResult
}

func (*fakeConn) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("fakeConn: Prepare is not supported")
}

func (*fakeConn) Close() error {
	return nil
}

func (*fakeConn) Begin() (driver.Tx, error) {
	return fakeTx{}, nil
}

func (c *fakeConn) QueryContext(_ context.Context, _ string, _ []driver.NamedValue) (driver.Rows, error) {
	if c.result.queryErr != nil {
		return nil, c.result.queryErr
	}
	return &fakeRows{row: c.result.row}, nil
}

type fakeTx struct{}

func (fakeTx) Commit() error   { return nil }
func (fakeTx) Rollback() error { return nil }

// fakeRows returns row once (if non-nil), then io.EOF — enough to model the
// single-row `pg_try_advisory_xact_lock` result set, including the no-rows case.
type fakeRows struct {
	row      []driver.Value
	returned bool
}

func (*fakeRows) Columns() []string {
	return []string{"pg_try_advisory_xact_lock"}
}

func (*fakeRows) Close() error {
	return nil
}

func (r *fakeRows) Next(dest []driver.Value) error {
	if r.row == nil || r.returned {
		return io.EOF
	}
	r.returned = true
	copy(dest, r.row)
	return nil
}

func TestTryTakeAdvisoryLock(t *testing.T) {
	errFakeQuery := errors.New("fake query error")

	t.Run("succeeds", func(t *testing.T) {
		t.Run("lock acquired", func(t *testing.T) {
			tx := beginFakeTx(t, fakeResult{row: []driver.Value{true}})
			err := tryTakeAdvisoryLock(context.Background(), tx, "test-key")
			assert.NoError(t, err)
		})
	})

	t.Run("fails", func(t *testing.T) {
		tests := []struct {
			name            string
			result          fakeResult
			wantErrIs       error
			wantErrContains string
		}{
			{
				name:      "query error",
				result:    fakeResult{queryErr: errFakeQuery},
				wantErrIs: errFakeQuery,
			},
			{
				name:            "no rows returned",
				result:          fakeResult{},
				wantErrContains: "no rows returned",
			},
			{
				name:   "scan error",
				result: fakeResult{row: []driver.Value{"not-a-bool"}},
			},
			{
				name:      "data lock taken",
				result:    fakeResult{row: []driver.Value{false}},
				wantErrIs: ErrDataLockTaken,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				tx := beginFakeTx(t, tt.result)
				err := tryTakeAdvisoryLock(context.Background(), tx, "test-key")
				require.Error(t, err)
				if tt.wantErrIs != nil {
					assert.ErrorIs(t, err, tt.wantErrIs)
				}
				if tt.wantErrContains != "" {
					assert.Contains(t, err.Error(), tt.wantErrContains)
				}
			})
		}
	})
}
