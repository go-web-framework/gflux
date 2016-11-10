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
	handlers map[string]http.Handler
	method string
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

func (m *Mux) handle(path string, mw []Middleware, h http.Handler, method string) *Route {
	r := &Route{
		path:       path,
		middleware: mw,
		handler:    h,
		method: method,
	}
	// insert returns an error if same r.path
	// with at least one of the same HTTP methods already exists.
	if err := m.trie.insert(r); err != nil {
		panic(err)
	}
	return r
}

func (m *Mux) Handle(path string, mw []Middleware, h http.Handler) *Route {
	return m.handle(path, mw, h, MethodAll)
}

func (m *Mux) GET(path string, mw []Middleware, h http.Handler) *Route {
	return m.handle(path, mw, h, MethodGet)
}

func (m *Mux) POST(path string, mw []Middleware, h http.Handler) *Route {
	return m.handle(path, mw, h, MethodPost)
}

func (m *Mux) PUT(path string, mw []Middleware, h http.Handler) *Route {
	return m.handle(path, mw, h,  MethodPut)
}

func (m *Mux) PATCH(path string, mw []Middleware, h http.Handler) *Route {
	return m.handle(path, mw, h, MethodPatch)
}

func (m *Mux) DELETE(path string, mw []Middleware, h http.Handler) *Route {
	return m.handle(path, mw, h, MethodDelete)
}

func (m *Mux) HEAD(path string, mw []Middleware, h http.Handler) *Route {
	return m.handle(path, mw, h, MethodHead)
}

func (m *Mux) OPTIONS(path string, mw []Middleware, h http.Handler) *Route {
	return m.handle(path, mw, h, MethodOptions)
}

type Params map[string]string

func GetParams(r *http.Request) Params {
	return r.Context().Value(paramsCtxKey).(Params)
}

func SetParams(r *http.Request, p Params) *http.Request {
	c := r.Context()
	c = context.WithValue(c, paramsCtxKey, p)
	r = r.WithContext(c) 
	return r
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

	if !found || route == nil{
		if  m.NotFound != nil {
			m.NotFound.ServeHTTP(w, r)
		} else {
			http.NotFound(w, r)

		}
		return
	}

	if r.Method == MethodOptions {
		var allowed string
		for key, _ := range route.handlers {
       if len(allowed) == 0 {
				allowed = key
			} else {
				allowed += ", " + key
			}
    }
    w.Header().Set("Allow", allowed)
		return
	}

	handl, err := route.getAllowed(r.Method)
	if err != nil {
		http.Error(w,
					http.StatusText(http.StatusMethodNotAllowed),
					http.StatusMethodNotAllowed,
				)
		return
	}
	
	r = SetParams(r, Params(p))
	run(w, r, route.middleware, handl)
	
	return
}

func (m *Mux) SetNotFound(handler http.Handler) {
 		m.NotFound = handler
 }


func (r *Route) getAllowed(method string) (http.Handler, error){ 
 	if handl, ok := r.handlers[MethodAll]; ok {
		return handl, nil
	}

	if handl, ok := r.handlers[method]; ok {
		return handl, nil
	}

	return nil, errors.New("Not allowed")
}
