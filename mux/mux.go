package mux

import (
	"errors"
	"net/http"
)

type key int

const paramsCtxKey key = 0

var Stop = errors.New("exit request handling")

// Middleware represents an HTTP middleware function.
type Middleware func(w http.ResponseWriter, r *http.Request) error

type Route struct {
	Path       string
	Middleware []Middleware
	Handler    http.Handler

	// Methods is the list of allowed HTTP methods.
	// If len(Methods) == 0, all HTTP methods are allowed.
	Methods []string
}

type Mux struct {
	trie     *Trie
	notFound http.Handler
}

func New() *Mux {
	return &Mux{
		trie: NewTrie(),
	}
}

func (m *Mux) Handle(path string, mw []Middleware, h http.Handler) *Route {
	r := &Route{
		Path:       path,
		Middleware: mw,
		Handler:    h,
	}
	if err := m.trie.insert(r); err != nil {
		panic(err)
	}
	return r
}

func (m *Mux) GET(path string, mw []Middleware, h http.Handler) *Route {
	return m.Handle(path, mw, h).Allow("GET")
}

func (m *Mux) POST(path string, mw []Middleware, h http.Handler) *Route {
	return m.Handle(path, mw, h).Allow("POST")
}

func (m *Mux) PUT(path string, mw []Middleware, h http.Handler) *Route {
	return m.Handle(path, mw, h).Allow("PUT")
}

func (m *Mux) DELETE(path string, mw []Middleware, h http.Handler) *Route {
	return m.Handle(path, mw, h).Allow("DELETE")
}

func (m *Mux) HEAD(path string, mw []Middleware, h http.Handler) *Route {
	return m.Handle(path, mw, h).Allow("HEAD")
}

func (r *Route) Allow(methods ...string) *Route {
	r.Methods = append(r.Methods, difference(r.Methods, methods)...)
	return r
}

// difference returns the strings in a that aren't in b.
func difference(a, b []string) []string {
	var ret []string

	m := make(map[string]bool, len(b))
	for _, s := range b {
		m[s] = true
	}

	for _, s := range a {
		if m[s] {
			continue
		}
		ret = append(ret, s)
	}

	return ret
}

func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO

	e, _, found := m.trie.Get(r.URL.Path)
	if !found {
		m.HandleNotFound(w, r)
		return
	}

	for _, middleware := range e.Middleware {
		err := middleware(w, r)
		if err != nil { //err == Stop
			//middleware error
		}
	}

	//varValues := map[string]string{"Var1": "aaa", "2": "2"}
	//varValues := "aaa"
	//ctx := r.Context()
	//ctx = context.WithValue(ctx, varKey, varValues)
	//r = r.WithContext(ctx)

	//method check?
	e.Handler.ServeHTTP(w, r)
}

// NotFound the mux custom 404 handler
func (m *Mux) SetNotFound(handler http.Handler) {
	m.notFound = handler
}

// HandleNotFound handle when a request does not match a registered handler.
func (m *Mux) HandleNotFound(rw http.ResponseWriter, req *http.Request) {
	if m.notFound != nil {
		m.notFound.ServeHTTP(rw, req)
	} else {
		http.NotFound(rw, req)
	}
}

/*
//Vars
func Vars(r *http.Request) string{
	return r.Context().Value(varKey).(string)
	//return r.Context().Value(varKey).(map[string]string)
}
*/
