package mux

import (
	"errors"
	"net/http"
	//"context"
)

const varKey int = 0

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
	r, err := m.radix.NewRoute(path, h, mw, method)
	if err != nil {
		return nil
	} else {
		return r
	}
}

func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request){
	//e, _ := m.tree.Get(r.URL.Path);
	e, found, _ := m.radix.Get(r.URL.Path);
	if !found {
			m.HandleNotFound(w, r)
			return
	}

	for _, middleware := range e.Middleware{
		err := middleware(w, r)
		if (err != nil){ //err == Stop
  		//middleware error
		}
	}
	//currently arbitrary values
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

func (m *Mux) AllowMethod(path string, method ...string){
		m.radix.UpdateRouteMethods(path, method...)
}

/*
//Vars
func Vars(r *http.Request) string{
	return r.Context().Value(varKey).(string)
	//return r.Context().Value(varKey).(map[string]string)
}
*/

