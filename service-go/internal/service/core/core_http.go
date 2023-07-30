package core

import (
	"context"
	"errors"
	"math/rand"
	"net/http"
	"svc/internal/z"
	"time"

	"github.com/julienschmidt/httprouter"
)

var (
	ErrInvalidHTTPServer = errors.New("svc/core: invalid http server")
)

type h = http.Handler // type alias for shortcut

func (x *Core) buildHTTP(ctx context.Context) (*Core, error) {
	_, trc, _ := z.OTel(ctx, "", "", "")
	ctx, span := trc.Start(ctx, "internal/service/core:Core.buildHTTP")
	defer span.End()

	if !x.c.ServingHTTP {
		return x, nil
	}

	srv := &http.Server{}
	x.c.WithHTTPServer(srv)
	if srv == nil {
		return nil, ErrInvalidHTTPServer
	}

	mux := &httprouter.Router{}
	x.c.WithHTTPRouter(mux)
	if mux == nil {
		return nil, ErrInvalidHTTPServer
	}

	srv.Handler = x.registerHTTP(ctx, mux)
	x.httpserver = srv
	return x, nil
}

func (x *Core) startHTTP(ctx context.Context) error {
	_, trc, _ := z.OTel(ctx, "", "", "")
	ctx, span := trc.Start(ctx, "internal/service/core:Core.startHTTP")
	defer span.End()

	if !x.c.ServingHTTP {
		return nil
	}
	if x.httpserver.TLSConfig != nil {
		return x.httpserver.ListenAndServeTLS("", "")
	} else {
		return x.httpserver.ListenAndServe()
	}
}

func (x *Core) stopHTTP(ctx context.Context) error {
	_, trc, _ := z.OTel(ctx, "", "", "")
	ctx, span := trc.Start(ctx, "internal/service/core:Core.stopHTTP")
	defer span.End()

	if !x.c.ServingHTTP {
		return nil
	}
	return x.httpserver.Shutdown(ctx)
}

func (x *Core) registerHTTP(ctx context.Context, mux *httprouter.Router) *httprouter.Router {
	_, trc, _ := z.OTel(ctx, "", "", "")
	ctx, span := trc.Start(ctx, "internal/service/core:Core.registerHTTP")
	defer span.End()

	mux.PanicHandler = z.When(mux.PanicHandler == nil, panicHandler{x}.Fn, nil)
	mux.NotFound = z.When(mux.NotFound == nil, h(notFoundHandler{x}), nil)
	mux.Handler(http.MethodGet, "/wait", waitHandler{x})
	return mux
}

// ---------------------------------------------------------------------------------------------------------------------
type panicHandler struct{ *Core }

func (x panicHandler) Fn(w http.ResponseWriter, r *http.Request, i interface{}) {
	ctx := r.Context()
	_, trc, _ := z.OTel(ctx, "", "", "")
	ctx, span := trc.Start(ctx, "internal/service/core:panicHandler.Fn")
	defer span.End()

	code := http.StatusInternalServerError
	http.Error(w, http.StatusText(code), code)
}

// ---------------------------------------------------------------------------------------------------------------------
type notFoundHandler struct{ *Core }

func (x notFoundHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	_, trc, _ := z.OTel(ctx, "", "", "")
	ctx, span := trc.Start(ctx, "internal/service/core:notFoundHandler.ServeHTTP")
	defer span.End()
	//
}

// ---------------------------------------------------------------------------------------------------------------------
type waitHandler struct{ *Core }

func (x waitHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	_, trc, _ := z.OTel(ctx, "", "", "")
	ctx, span := trc.Start(ctx, "internal/service/core:waitHandler.ServeHTTP")
	defer span.End()

	n := 1 + rand.Intn(3_000)
	<-time.After(time.Duration(n) * time.Millisecond)
}
