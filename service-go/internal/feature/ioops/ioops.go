package ioops

import (
	"context"
	"svc/internal/repository/sqlcore"
)

type Configuration struct{}

type Dependency struct {
	SQL_Core
}

func New(ctx context.Context, c Configuration, d Dependency) (*IO_Ops, error) {
	return &IO_Ops{c, d}, nil
}

type IO_Ops struct {
	c Configuration
	d Dependency
}

// ---------------------------------------------------------------------------------------------------------------------

var _ SQL_Core = (*sqlcore.SQL_Core)(nil)

type SQL_Core interface {
	DummyUpsert(ctx context.Context, req sqlcore.DummyUpsertRequest) (res sqlcore.DummyUpsertResponse, err error)
}
