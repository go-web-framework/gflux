// package api allows for easy creation of REST APIs
package api

import (
	"../mux"
	"fmt"
	"net/http"
)

type API struct {
	DB  *Orm
	mux *mux.Mux
}

func New(dbDriver string, dbPath string) *API {
	db := InitDB(dbDriver, dbPath)
	return &API{DB: db, mux: mux.New()}
}

func (a *API) Close() {
	a.DB.Close()
}

// create a new resource for the api
func (a *API) NewResource(name string, structType interface{}) *Resource {
	// create resource
	res := NewResource(name, structType, a)

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
