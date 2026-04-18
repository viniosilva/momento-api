package nethttp

import (
	"net/http"
	"time"

	nethttp_port "momento/pkg/nethttp/port"
	nethttp_recovery "momento/pkg/nethttp/recovery"
	nethttp_requestid "momento/pkg/nethttp/requestid"
	nethttp_timeout "momento/pkg/nethttp/timeout"
	nethttp_utils "momento/pkg/nethttp/utils"
)

const defaultTimeout = 10 * time.Second

type Chain struct {
	middlewares []nethttp_utils.MiddlewareFunc
	timeout     *time.Duration
}

type ChainOption func(*Chain)

func WithTimeout(timeout *time.Duration) ChainOption {
	return func(c *Chain) {
		c.timeout = timeout
	}
}

func NewChain(middlewares ...nethttp_utils.MiddlewareFunc) *Chain {
	return &Chain{middlewares: middlewares}
}

func NewDefaultChain(logger nethttp_port.LoggerSlog, opts ...ChainOption) *Chain {
	c := NewChain(
		nethttp_recovery.RecoveryMiddleware(logger),
		nethttp_requestid.RequestIDMiddleware,
	)

	for _, opt := range opts {
		opt(c)
	}

	if c.timeout == nil {
		timeout := defaultTimeout
		c.timeout = &timeout
	}
	c.AddMiddleware(nethttp_timeout.TimeoutMiddleware(*c.timeout))

	return c
}

func (c *Chain) AddMiddleware(middleware nethttp_utils.MiddlewareFunc) {
	c.middlewares = append(c.middlewares, middleware)
}

func (c *Chain) Then(h http.Handler) http.Handler {
	for i := len(c.middlewares) - 1; i >= 0; i-- {
		h = c.middlewares[i](h)
	}

	return h
}

func (c *Chain) ThenFunc(fn http.HandlerFunc) http.Handler {
	return c.Then(fn)
}
