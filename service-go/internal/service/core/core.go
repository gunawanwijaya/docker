package core

import (
	"context"
	"errors"
	"net/http"
	"svc/internal/feature/ioops"
	"svc/internal/service"
	"svc/internal/z"

	"github.com/julienschmidt/httprouter"
)

type Configuration struct {
	ServingHTTP bool

	WithHTTPServer func(srv *http.Server)       `json:"-" yaml:"-"`
	WithHTTPRouter func(mux *httprouter.Router) `json:"-" yaml:"-"`
}

type Dependency struct {
	IO_Ops
}

func New(ctx context.Context, c Configuration, d Dependency) (service.Service, error) {
	log, trc, _ := z.OTel(ctx, "", "", "")
	ctx, span := trc.Start(ctx, "internal/service/core:_.New")
	defer span.End()

	log.Debug().Stringer("debugKey", debugKey).Send()

	x := &Core{c, d, nil}
	x, err := x.buildHTTP(ctx)
	span.RecordError(err)

	return x, err
}

type Core struct {
	c Configuration
	d Dependency

	httpserver *http.Server
}

func (x *Core) Start(ctx context.Context) error {
	_, trc, _ := z.OTel(ctx, "", "", "")
	ctx, span := trc.Start(ctx, "internal/service/core:Core.Start")
	defer span.End()

	chErr := make(chan error)
	switch {
	case x.c.ServingHTTP:
		go func() { chErr <- x.startHTTP(ctx) }()
	}
	return <-chErr
}

func (x *Core) Stop(ctx context.Context) error {
	_, trc, _ := z.OTel(ctx, "", "", "")
	ctx, span := trc.Start(ctx, "internal/service/core:Core.Stop")
	defer span.End()

	return errors.Join(
		x.stopHTTP(ctx),
	)
}

var (
	debugKey, _ = z.Strategy.Encoding.Base64(1).Encode(
		z.Strategy.Keygen.Nonce(64),
	)
)

// ---------------------------------------------------------------------------------------------------------------------

var _ IO_Ops = (*ioops.IO_Ops)(nil)

type IO_Ops interface {
	OnDatabase(ctx context.Context, req ioops.OnDatabaseRequest) (res ioops.OnDatabaseResponse, err error)
}
