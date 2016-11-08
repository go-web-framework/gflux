// package api allows for easy creation of REST APIs
package api

import (
	"github.com/go-web-framework/gflux/mux"
	"net/http"
	"reflect"
	"fmt"
)

type Resource struct {
	Name     string
	Type     reflect.Type
	api      *API
	ItemHandlers map[string]func(interface{}, http.ResponseWriter, []string)
	CollectionHandlers map[string]func([]interface{}, http.ResponseWriter, []string)
}

func NewResource(name string, structType interface{}, api *API) *Resource {
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
			obj = api.DB.FindAndDelete(res.Type, res.Name, id)
		} else {
			obj = api.DB.Find(res.Type, res.Name, id)
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
