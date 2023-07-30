package sqlcore

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"svc/internal/z"
	"svc/internal/z/zsql"
)

type DummySelectRequest struct {
	List []DummySelectRequest

	Column string
}

type DummySelectResponse struct {
	List []DummySelectResponse

	Column string
}

func (x *SQL_Core) DummySelect(ctx context.Context, req DummySelectRequest) (res DummySelectResponse, err error) {
	q := z.Yield(dummy_select)() // make dummy_select immutable

	err = zsql.Wrap.Query(x.d.Conn.QueryContext(ctx, q)).Scan(func(i int) (pointers []any) {
		res.List = append(res.List, DummySelectResponse{})
		return []any{
			&res.List[i].Column,
		}
	})
	return
}

type DummyUpsertRequest struct {
	List []DummyUpsertRequest

	Column string
}

type DummyUpsertResponse struct {
	List []DummyUpsertResponse

	Column string
}

func (x *SQL_Core) DummyUpsert(ctx context.Context, req DummyUpsertRequest) (res DummyUpsertResponse, err error) {
	if len(req.List) == 0 {
		req.List = append(req.List, req)
	}
	l := req.List
	// -----------------------------------------------------------------------------------------------------------------
	// validate request
	m := make(map[string]struct{})
	for i := 0; i < len(l); i++ {
		// remove "" and duplicate from slice
		if _, found := m[l[i].Column]; found || l[i].Column == "" {
			copy(l[i:], l[i+1:])               // Shift a[i+1:] left one index.
			l[len(l)-1] = DummyUpsertRequest{} // Erase last element (write zero value).
			l = l[:len(l)-1]                   // Truncate slice.
			i = i - 1
			continue
		}
		m[l[i].Column] = struct{}{}
	}
	if len(l) < 1 {
		err = fmt.Errorf("aborting the execution due to no valid entry")
		return res, err
	}

	// -----------------------------------------------------------------------------------------------------------------
	// query manipulation so that multirows insert can be done in one go
	q := z.Yield(dummy_upsert)()      // make dummy_upsert immutable
	q = q[:strings.LastIndex(q, ";")] // because contains ; at the end of the query we want to remove it
	b := strings.Builder{}            // use strings.Builder
	b.Grow(len(q) + (5*len(l) - 1))   // grow at least 5 chars `,($X)` on each iter
	_, _ = b.WriteString(q)           // because contains ; at the end of the query we want to remove it
	a := make([]any, len(l))          // prepare the list of arguments to be passed
	for i := 0; i < len(l); i++ {     // iteration
		a[i] = l[i].Column // assign arguments to be passed
		if i > 0 {
			_, _ = b.WriteString(",($")
			_, _ = b.WriteString(strconv.Itoa(i + 1))
			_, _ = b.WriteString(")")
		}
	}
	_, _ = b.WriteString(";") // since we remove the ;, we add it back for sweetness
	q = b.String()            // reassign the new query
	// -----------------------------------------------------------------------------------------------------------------
	err = zsql.Wrap.Transaction(x.d.BeginTx(ctx, &sql.TxOptions{})).Do(func() error {
		rowsAffected := 0
		err = zsql.Wrap.Exec(x.d.Conn.ExecContext(ctx, q, a...)).Scan(&rowsAffected, nil)
		if err != nil {
			err = fmt.Errorf("%w: %s", err, q)
		}
		if len(l) != rowsAffected {
			err = fmt.Errorf("rolling back due to difference of expected input (%d) & actual record (%d)",
				len(l),
				rowsAffected,
			)
		}
		return err
	})
	if err != nil {
		return res, err
	}
	for _, v := range l {
		res.List = append(res.List, DummyUpsertResponse{
			Column: v.Column,
		})
	}
	return res, err
}
