package sqlcore_test

import (
	"context"
	"database/sql"
	"svc/internal/repository/sqlcore"
	"svc/internal/z/zsql"
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDummyUpsert(t *testing.T) {
	ctx := context.Background()
	dbTest(t, false, func(db *sql.DB) {
		func() {
			_, err := db.ExecContext(ctx, `
			CREATE TABLE IF NOT EXISTS tbl_dummy(_column TEXT);
			`)
			require.NoErrorf(t, err, "unable to seed: %s", err)
		}()

		sqlcoreRepository, err := sqlcore.New(ctx,
			sqlcore.Configuration{},
			sqlcore.Dependency{db},
		)
		assert.NoErrorf(t, err, "unable to initialize [sqlcore]: %s", err)

		{
			res, err := sqlcoreRepository.DummyUpsert(ctx, sqlcore.DummyUpsertRequest{
				List: []sqlcore.DummyUpsertRequest{
					{Column: "a"}, // INSERTED
					{Column: "b"}, // INSERTED
					{Column: "c"}, // INSERTED
					{Column: "a"}, // IGNORED: duplicate
					{Column: ""},  // IGNORED: empty
				},
			})
			assert.NoErrorf(t, err, "unable to request DummyUpsert [sqlcore]: %s", err)
			_ = res
		}
		{
			res, err := sqlcoreRepository.DummySelect(ctx, sqlcore.DummySelectRequest{})
			assert.NoErrorf(t, err, "unable to request DummySelect [sqlcore]: %s", err)
			assert.Len(t, res.List, 3) // from previous insert
		}
	})
}

// dbTest
//
// Here instead of using sqlmock, we opt to use a live postgres database to let us test for a deeper integration.
// The reason why we avoid sqlmock is due to the query is not properly compatible with the target database, e.g.
// calling a function that is vendor specific are guarantee not to return any error, unless we define in the sqlmock
//
// Note that the newly created database on the fn callback still empty, if we want to seed the database, do it in the
// callback. Additionally, we can also retain the database for debugging purposes.
func dbTest(t *testing.T, retain bool, fn func(db *sql.DB)) {
	if fn == nil {
		return
	}

	dsn := func(dbName string) string {
		return "postgres://postgres-mock:postgres-mock@localhost:5431/" + dbName + "?sslmode=disable"
	}

	// connect to the mock server default database
	dsnMock := dsn("postgres-mock")
	dbMock, err := zsql.Open(dsnMock)
	require.NoErrorf(t, err, "unable to connect to mock database: %s: %s", err, dsnMock)

	// create a new temp database owned by current user
	dbNameTemp := xid.New().String()
	_, err = dbMock.Exec("CREATE DATABASE " + dbNameTemp + " WITH OWNER \"postgres-mock\"")
	require.NoErrorf(t, err, "unable to create temp database: %s: %s", err, dbNameTemp)

	// close connection to the server
	err = dbMock.Close()
	require.NoErrorf(t, err, "unable to close mock database: %s", err)

	// connect to the mock server temp database
	dsnTest := dsn(dbNameTemp)
	dbTest, err := zsql.Open(dsnTest)
	require.NoErrorf(t, err, "unable to connect to test database: %s: %s", err, dsnTest)

	fn(dbTest) // <- the callback

	// test is done, close connection to the server
	err = dbTest.Close()
	require.NoErrorf(t, err, "unable to close test database: %s", err)

	if !retain {
		// connect to the mock server default database
		dbMock, err := zsql.Open(dsnMock)
		require.NoErrorf(t, err, "unable to connect to mock database: %s: %s", err, dsnMock)

		// drop temp database
		_, err = dbMock.Exec("DROP DATABASE " + dbNameTemp)
		require.NoErrorf(t, err, "unable to drop temp database: %s: %s", err, dbNameTemp)

		// close connection to the server
		err = dbMock.Close()
		require.NoErrorf(t, err, "unable to close mock database: %s", err)
	}
}
