package api

import (
	"../mux" // TODO: change to full path
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
)

// For inline http.Handler creation
type ProtoHttpHandler struct {
	ServeHTTPMethod func(w http.ResponseWriter, r *http.Request)
}

func (h ProtoHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.ServeHTTPMethod(w, r)
	return
}

type Model struct {
	gorm.Model
}

type API struct {
	db  *gorm.DB
	mux *mux.Mux
}

type Resource struct {
	Name string
}

func New(dbName string) (*API, error) {
	db, err := gorm.Open("sqlite3", dbName)
	if err != nil {
		return nil, err
	}
	return &API{db: db, mux: mux.New()}, nil
}

/////////  API  /////////////////////////////

func (a *API) Close() {
	a.db.Close()
}

func (a *API) NewResource(name string, i interface{}) *Resource {
	h := struct{ ProtoHttpHandler }{}
	h.ServeHTTPMethod = func(w http.ResponseWriter, r *http.Request) {
		id := mux.GetParams(r)["id"]
		fmt.Fprintf(w, "<h1>You've reached resource "+name+" with id "+id+"!</h1>")
		return
	}
	a.mux.Handle("/"+name+"/{id}", nil, h)
	return &Resource{name}
	// TODO: use i for database initialization
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
