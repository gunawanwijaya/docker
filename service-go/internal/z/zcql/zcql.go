package zcql

import (
	"context"
	"errors"

	"github.com/gocql/gocql"
)

func A() {
	c := gocql.NewCluster()
	localDC := ""

	// Enable token aware host selection policy, if using multi-dc cluster set a local DC.
	fallback := gocql.RoundRobinHostPolicy()
	if localDC != "" {
		fallback = gocql.DCAwareRoundRobinPolicy(localDC)
	}
	c.PoolConfig.HostSelectionPolicy = gocql.TokenAwareHostPolicy(fallback)

	// If using multi-dc cluster use the "local" consistency levels.
	if localDC != "" {
		c.Consistency = gocql.LocalQuorum
	}

	s, _ := c.CreateSession()
	q := s.Query("")
	defer q.Release()
	ctx, cancelCause := context.WithCancelCause(q.Context())
	cancelCause(errors.New("soemthing"))
	q = q.WithContext(ctx)
}
