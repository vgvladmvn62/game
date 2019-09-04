package server

import (
	"net/http"
	"os"

	"github.com/go-http-utils/logger"
	"github.com/justinas/alice"
)

// Log logs data to stdout with github.com/go-http-utils/logger.
func Log(t logger.Type) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return logger.Handler(h, os.Stdout, t)
	}
}

// HeaderHandler acts as middleware adding headers.
type HeaderHandler struct {
	key   string
	value string
	h     http.Handler
}

// ServeHTP serves HTML document.
func (h HeaderHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add(h.key, h.value)
	h.h.ServeHTTP(w, r)
}

// Middleware holds different chains of middleware and merges them as needed.
type Middleware struct {
	common alice.Chain
	api    alice.Chain
	html   alice.Chain
}

// API returns middleware chain containing common and api chain.
func (m *Middleware) API(h http.HandlerFunc) http.Handler {
	return m.common.Extend(m.api).ThenFunc(h)
}

// HTML returns middleware chain containing common and html chain.
func (m *Middleware) HTML(h http.HandlerFunc) http.Handler {
	return m.common.Extend(m.html).ThenFunc(h)
}

// Header creates HeaderHandler middleware.
func Header(key, value string) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return HeaderHandler{
			key:   key,
			value: value,
			h:     handler,
		}
	}
}
