package api

import (
	"../mux"
	"net/http"
	"reflect"
	"fmt"
	"strings"
)

// TODO: Define this for the godoc.
type Resource struct {
	Name     string
	Type     reflect.Type
	api      *API
	ItemHandlers map[string]func(interface{}, http.ResponseWriter, []string)
	CollectionHandlers map[string]func([]interface{}, http.ResponseWriter, []string)
}

// SetItemHandler overrides or adds the given handler to be called
// when a request with the specified method is made to /name/{id}.
// The function handler receives the item with the requested id.
// If DELETE, the object will have already been deleted.
// If no object exists in the database with the id, the object passed
// to the handler function is nil.
// The handler function also receives a slice containing the accepts in
// the form of ["application/json", "application/xml"].
func (res *Resource) SetItemHandler(method string, handler func(object interface{}, w http.ResponseWriter, accepts []string)) {
	res.ItemHandlers[strings.ToUpper(method)] = handler
}

// SetItemHandler overrides or adds the given handler to be called
// when a request with the specified method is made to /name/{id}.
// The function handler receives all items in the database.
// If the database is empty, it receives nil.
// The handler function also receives a slice containing the accepts in
// the form of ["application/json", "application/xml"].
func (res *Resource) SetCollectionHandler(method string, handler func(objects []interface{}, w http.ResponseWriter, accepts []string)) {
	res.CollectionHandlers[strings.ToUpper(method)] = handler
}

func newResource(name string, structType interface{}, api *API) *Resource {
	t := reflect.TypeOf(structType)

	// If ptr, dereference schema type
	if t.Kind().String() == "ptr" {
		t = t.Elem()
	}

	r := Resource{Name: name, Type: t, api: api}
	
	r.ItemHandlers = make(map[string]func(interface{}, http.ResponseWriter, []string))
	r.ItemHandlers["GET"] = defaultItemGET
	r.ItemHandlers["DELETE"] = defaultItemDELETE
	
	r.CollectionHandlers = make(map[string]func([]interface{}, http.ResponseWriter, []string))
	r.CollectionHandlers["GET"] = defaultCollectionGET

	return &r
}

type itemHandler struct {
	res *Resource
}

type collectionHandler struct {
	res *Resource
}

func (h itemHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res := h.res
	api := res.api

	// find id from request
	id := mux.GetParams(r)["id"]
	
	// Check if handler has been implemented for the request method
	_, exists := res.ItemHandlers[r.Method]
	if exists == true {
		// query database
		var obj interface{}
		if(r.Method == "DELETE"){
			obj = api.DB.DeleteById(res.Type, res.Name, id)
		} else {
			obj = api.DB.FindById(res.Type, res.Name, id)
		}
		
		// call handler
		res.ItemHandlers[r.Method](obj, w, []string{"application/json"})
	} else {
		fmt.Println(r.RequestURI + " does not have a " + r.Method + " method defined")
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
	}

	return
}

func (h collectionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res := h.res
	api := res.api
	
	// Check if handler has been implemented for the request method
	_, exists := res.CollectionHandlers[r.Method]
	if exists == true {
		// read from database
		objs := api.DB.FindAll(res.Type, res.Name)
		res.CollectionHandlers[r.Method](objs, w, []string{"application/json"})
	} else {
		fmt.Println(r.RequestURI + " does not have a " + r.Method + " method defined")
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
	}

	return
}
