// package api allows for easy creation of REST APIs
package api

import (
	"../mux"
	"fmt"
	"net/http"
)

// TODO: define this for the godoc
type API struct {
	DB  *Orm
	mux *mux.Mux
}

// New creates a new API with the specified database connection.
func New(dbDriver string, dbPath string) *API {
	db := InitDB(dbDriver, dbPath)
	return &API{DB: db, mux: mux.New()}
}

// Close closes the API's database connection.
func (a *API) Close() {
	a.DB.Close()
}

// NewResource creates a new resource with a new database table with the given name
// and routes to /name/{id} and to /name/.
// The new resource has GET, DELETE and POST handlers defined.
func (a *API) NewResource(name string, structType interface{}) *Resource {
	// create resource
	res := newResource(name, structType, a)

	// create database table
	a.DB.CreateTable(name, structType)

	// create handlers
	hItem := itemHandler{res}
	hCollection := collectionHandler{res}

	// assign handler
	a.mux.Handle("/"+name+"/{id}", nil, hItem)
	a.mux.Handle("/"+name, nil, hCollection)

	return res
}

// Serve serves the API to localhost at the specified port.
// The default port is :8080.
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
