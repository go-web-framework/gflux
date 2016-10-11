package mux

import (
	"context"
	"errors"
	"net/http"
)

type key int

const paramsCtxKey key = 0

var Stop = errors.New("exit request handling")

// Middleware represents an HTTP middleware function.
type Middleware func(w http.ResponseWriter, r *http.Request) error

type Route struct {
	path       string
	middleware []Middleware
	handler    http.Handler

	// methods is the list of allowed HTTP methods.
	// If len(Methods) == 0, all HTTP methods are allowed.
	methods []string
}

type Mux struct {
	trie     *Trie
	NotFound http.Handler
}

func New() *Mux {
	return &Mux{
		trie: NewTrie(),
	}
}

func (m *Mux) handle(path string, mw []Middleware, h http.Handler, methods []string) *Route {
	r := &Route{
		path:       path,
		middleware: mw,
		handler:    h,
		methods:    methods,
	}
	// TODO: insert currently returns an error if r.path already exists.
	// Instead, it should return an error only if same r.path
	// with at least one of the same HTTP methods already exists.
	if err := m.trie.insert(r); err != nil {
		panic(err)
	}
	return r
}

func (m *Mux) Handle(path string, mw []Middleware, h http.Handler) *Route {
	return m.handle(path, mw, h, nil)
}

func (m *Mux) GET(path string, mw []Middleware, h http.Handler) *Route {
	return m.handle(path, mw, h, []string{"GET"})
}

func (m *Mux) POST(path string, mw []Middleware, h http.Handler) *Route {
	return m.handle(path, mw, h, []string{"POST"})
}

func (m *Mux) PUT(path string, mw []Middleware, h http.Handler) *Route {
	return m.handle(path, mw, h, []string{"PUT"})
}

func (m *Mux) PATCH(path string, mw []Middleware, h http.Handler) *Route {
	return m.handle(path, mw, h, []string{"PATCH"})
}

func (m *Mux) DELETE(path string, mw []Middleware, h http.Handler) *Route {
	return m.handle(path, mw, h, []string{"DELETE"})
}

func (m *Mux) HEAD(path string, mw []Middleware, h http.Handler) *Route {
	return m.handle(path, mw, h, []string{"HEAD"})
}

func (m *Mux) OPTIONS(path string, mw []Middleware, h http.Handler) *Route {
	return m.handle(path, mw, h, []string{"OPTIONS"})
}

type Params map[string]string

func GetParams(c context.Context) Params {
	return c.Value(paramsCtxKey).(Params)
}

func SetParams(c context.Context, p Params) {
	c = context.WithValue(c, paramsCtxKey, p)
}

func run(w http.ResponseWriter, r *http.Request, mw []Middleware, h http.Handler) {
	for _, m := range mw {
		if m != nil {
			if m(w, r) == Stop {
				return
			}
		}
	}
	if h != nil {
		h.ServeHTTP(w, r)
	}
}

func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route, p, found := m.trie.Get(r.URL.Path)

	if found {
		SetParams(r.Context(), Params(p))
		run(w, r, route.middleware, route.handler)
		return
	}

	if m.NotFound != nil {
		m.NotFound.ServeHTTP(w, r)
	} else {
		http.NotFound(w, r)
	}
}
