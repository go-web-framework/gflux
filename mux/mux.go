package mux

import (
	"errors"
	"net/http"
	)

var Stop = errors.New("exit request handling")

// Middleware represents an HTTP middlware function.
type Middleware func(w http.ResponseWriter, r *http.Request) error

type tree interface {
	Get(p string) (*Entry, bool)
	Set(p string, m *Entry)
}

type treeImpl map[string]*Entry

func (t treeImpl) Get(p string) (e *Entry, ok bool) {
	e, ok = t[p]
	return
}

func (t treeImpl) Set(p string, m *Entry) {
	t[p] = m
}

type Entry struct {
	Path       string
	Middleware []Middleware
	Handler    http.Handler
	Methods    []string // Allowed HTTP methods.
}

func (e *Entry) Allow(methods ...string) *Entry {
	e.Methods = append(e.Methods, methods...)
	return e
}

// Mux is a serve mux.
type Mux struct {
	tree tree
}

func New() *Mux {
	return &Mux{
		tree: make(treeImpl),
	}
}

func (m *Mux) Handle(path string, mw []Middleware, h http.Handler) *Entry {
	e := &Entry{
		Path:       path,
		Middleware: mw,
		Handler:    h,
		Methods:    nil,
	}
	m.tree.Set(path, e)
	return e
}

func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request){
	e, _ := m.tree.Get(r.URL.Path);
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
