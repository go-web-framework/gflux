package mux

import (
	"errors"
	"net/http"
)

var Stop = errors.New("exit request handling")

// Middleware represents an HTTP middlware function.
type Middleware func(w http.ResponseWriter, r *http.Request) error

// Mux is a serve mux.
type Mux struct {
	radix *Trie
	notFound http.Handler
}

func New() *Mux {
	return &Mux{
		radix: NewTrie(),
	}
}

func (m *Mux) Handle(path string, mw []Middleware, h http.Handler) *Route{
	method := []string{"Get"}
	return m.radix.NewRoute(path, h, mw, method)
}

func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request){
	//e, _ := m.tree.Get(r.URL.Path);
	e, found := m.radix.Get(r.URL.Path);
	if !found {
			m.HandleNotFound(w, r)
			return
	}

	//iterator best way?
	for _, middleware := range e.Middleware{
		err := middleware(w, r)
		if (err != nil){ //err == Stop
  		//middleware error
		}
	}


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

func (m *Mux) AllowMethod(path string, method ...string){
		m.radix.UpdateRouteMethods(path, method...)
}

