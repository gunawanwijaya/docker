package ioops

import "context"

type OnDatabaseRequest struct{}
type OnDatabaseResponse struct{}

func (x *IO_Ops) OnDatabase(ctx context.Context, req OnDatabaseRequest) (res OnDatabaseResponse, err error) {
	return
}
