package sqlcore

import (
	"context"
	"embed"
	"svc/internal/z/zsql"
)

type Configuration struct{}
type Dependency struct {
	zsql.Conn
}

func New(ctx context.Context, c Configuration, d Dependency) (*SQL_Core, error) {
	return &SQL_Core{c, d}, nil
}

type SQL_Core struct {
	c Configuration
	d Dependency
}

// ---------------------------------------------------------------------------------------------------------------------
var (
	_ embed.FS

	//go:embed query/dummy_upsert.sql
	dummy_upsert string

	//go:embed query/dummy_select.sql
	dummy_select string
)
