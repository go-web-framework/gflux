package api

import (
	"../mux" // TODO: change to full path
	"fmt"
	"net/http"
	"reflect"
)

// For inline http.Handler creation
type ProtoHttpHandler struct {
	ServeHTTPMethod func(w http.ResponseWriter, r *http.Request)
}

func (h ProtoHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.ServeHTTPMethod(w, r)
	return
}

type API struct {
	db  *Orm
	mux *mux.Mux
}

func New(dbPath string) (*API) {
	db := InitDB(dbPath)
	return &API{db: db, mux: mux.New()}
}

/////////  API  /////////////////////////////

func (a *API) Close() {
	a.db.Close()
}

// create a new resource for the api
func (a *API) NewResource(name string, i interface{}) *Resource {
	// create database table
	a.db.CreateTable(name, i)

	// create handler
	h := struct{ ProtoHttpHandler }{}
	h.ServeHTTPMethod = func(w http.ResponseWriter, r *http.Request) {
		id := mux.GetParams(r)["id"]
		fmt.Fprintf(w, "<h1>You've reached resource "+name+" with id "+id+"!</h1>")
		return
	}
	
	// assign handler
	a.mux.Handle("/"+name+"/{id}", nil, h)
	
	// create resource
	r := NewResource(name, reflect.TypeOf(i))
	
	return r
}

func (a *API) Serve(port ...string) {
	var portNum string
	if len(port) == 0 {
		portNum = ":8080"
	} else {
		portNum = port[0]
	}
	fmt.Println("Listening on " + portNum)
	http.ListenAndServe(portNum, a.mux)
}
