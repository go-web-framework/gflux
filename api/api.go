package api

import (
	"../mux" // TODO: change to full path
	"fmt"
	"net/http"
)

type API struct {
	db  *Orm
	mux *mux.Mux
}

func New(dbPath string) *API {
	db := InitDB(dbPath)
	return &API{db: db, mux: mux.New()}
}

/////////  API  /////////////////////////////

func (a *API) Close() {
	a.db.Close()
}

// create a new resource for the api
func (a *API) NewResource(name string, structType interface{}) *Resource {
	// create resource
	res := NewResource(name, structType, a)

	// create database table
	a.db.CreateTable(name, structType)

	// create handlers
	h := DefaultGETHandler{res}

	// assign handler
	a.mux.Handle("/"+name+"/{id}", nil, h)

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
