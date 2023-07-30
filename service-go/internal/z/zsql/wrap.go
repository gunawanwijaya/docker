package zsql

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"
)

var (
	ErrInvalidArguments   = errors.New("zsql: invalid arguments for scan")
	ErrInvalidTransaction = errors.New("zsql: invalid transaction")
	ErrNoColumns          = errors.New("zsql: no columns returned")
)

type MismatchColumnsError struct{ Col, Dst int }

func (err *MismatchColumnsError) Error() string {
	return fmt.Sprintf("zsql: mismatch %d columns on %d destinations",
		err.Col,
		err.Dst,
	)
}

// =============================================================================

//nolint:gochecknoglobals
var Wrap wrap

type wrap struct{}

// Exec will wrap `ExecContext` so that we can Scan later
//
//	Exec(cmd.ExecContext(ctx, "..."))
func (wrap) Exec(val sql.Result, err error) WrapExec {
	return wrapExec{val, err}
}

type WrapExec interface {
	// Scan the result of ExecContext that usually return numbers of rowsAffected and lastInsertID.
	// Keep in mind that some drivers are not support rowsAffected or lastInsertID
	Scan(rowsAffected *int, lastInsertID *int) error
}

type wrapExec struct {
	res sql.Result
	err error
}

func (x wrapExec) Scan(rowsAffected *int, lastInsertID *int) error {
	err := x.err
	if err != nil {
		return err
	}

	if x.res == nil {
		return ErrInvalidArguments
	}

	if rowsAffected != nil {
		// TODO: how to channel the error?
		n, _ := x.res.RowsAffected()
		*rowsAffected = int(n)
	}

	if lastInsertID != nil {
		// TODO: how to channel the error?
		n, _ := x.res.LastInsertId()
		*lastInsertID = int(n)
	}

	return err
}

// =============================================================================

// Query will wrap `QueryContext` so that we can Scan later
//
//	Query(cmd.QueryContext(ctx, "..."))
func (wrap) Query(val *sql.Rows, err error) WrapQuery {
	return wrapQuery{val, err}
}

type WrapQuery interface {
	// Scan accept do, a func that accept `i int` as index and returning list
	// of pointers.
	//  pointers == nil   // break the loop
	//  len(pointers) < 1 // skip the current loop
	//  len(pointers) > 0 // assign the pointer, MUST be same as the length of columns
	Scan(row func(i int) (pointers []any)) error
}

type wrapQuery struct {
	res *sql.Rows
	err error
}

func (x wrapQuery) Scan(row func(i int) []any) error {
	err := x.err
	if err != nil {
		return err
	} else if x.res == nil {
		return sql.ErrNoRows
	} else if err = x.res.Err(); err != nil {
		return err
	}
	defer x.res.Close()

	cols, err := x.res.Columns()
	if err != nil {
		return err
	} else if len(cols) < 1 {
		return ErrNoColumns
	}

	for i := 0; x.res.Next(); i++ {
		err = x.res.Err()
		if err != nil {
			return err
		}

		dest := row(i)
		if dest == nil { // nil dest
			break
		} else if len(dest) < 1 { // empty dest
			continue
		} else if len(dest) != len(cols) { // diff dest & cols
			return &MismatchColumnsError{len(cols), len(dest)}
		}

		err = x.res.Scan(dest...) // scan into pointers
		if err != nil {
			return err
		}
	}

	return err
}

// =============================================================================

// Query will wrap `QueryContext` so that we can Scan later
//
//	Query(cmd.QueryContext(ctx, "..."))
func (wrap) QueryRow(val *sql.Row, err error) WrapQueryRow {
	return wrapQueryRow{val, err}
}

type WrapQueryRow interface {
	Scan(dest ...any) error
	Err() error
}

type wrapQueryRow struct {
	res *sql.Row
	err error
}

func (x wrapQueryRow) Scan(dest ...any) error {
	return x.res.Scan(dest...)
}

func (x wrapQueryRow) Err() error {
	if x.err == nil {
		x.err = x.res.Err()
	}

	return x.err
}

// =============================================================================

// Transaction will wrap `Begin` so that we can Wrap later
//
//	Transaction(db.BeginTx(ctx, ...))
//
// Wrap the transaction and ends it with either COMMIT or ROLLBACK.
func (wrap) Transaction(tx *sql.Tx, err error) WrapTransaction {
	return &wrapTransaction{new(sync.Once), tx, err}
}

type WrapTransaction interface {
	// Do the transaction and ends it with either COMMIT or ROLLBACK
	//
	// non-nil error will ROLLBACK the transaction, but nil error will COMMIT the transaction
	Do(tx func() error) error
}

type wrapTransaction struct {
	once *sync.Once
	res  *sql.Tx
	err  error
}

func (x wrapTransaction) Do(cb func() error) error {
	if x.err != nil {
		return x.err
	}

	x.err = ErrInvalidTransaction
	x.once.Do(func() {
		if err := cb(); err != nil {
			if x.err = x.res.Rollback(); x.err == nil {
				x.err = err
			}
		} else {
			x.err = x.res.Commit()
		}
	})
	return x.err
}
