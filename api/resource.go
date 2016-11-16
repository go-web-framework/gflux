package api

import (
	"../mux"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
)

// TODO: Define this for the godoc.
type Resource struct {
	Name               string
	Type               reflect.Type
	api                *API
	ItemHandlers       map[string]func(http.ResponseWriter, *http.Request, interface{})
	CollectionHandlers map[string]func(http.ResponseWriter, *http.Request, []interface{})
	methods            map[string]struct{}
}

// SetItemHandler overrides or adds the given handler to be called
// when a request with the specified method is made to /name/{id}.
// The function handler receives the item with the requested id.
// If DELETE, the object will have already been deleted.
// If no object exists in the database with the id, the object passed
// to the handler function is nil.
func (res *Resource) SetItemHandler(method string, handler func(w http.ResponseWriter, r *http.Request, object interface{})) *Resource {
	res.ItemHandlers[strings.ToUpper(method)] = handler
	return res
}

// SetItemHandler overrides or adds the given handler to be called
// when a request with the specified method is made to /name/{id}.
// The function handler receives all items in the database.
// If the database is empty, it receives nil.
func (res *Resource) SetCollectionHandler(method string, handler func(w http.ResponseWriter, r *http.Request, object []interface{})) *Resource {
	res.CollectionHandlers[strings.ToUpper(method)] = handler
	return res
}

//TODO: godoc
func (res *Resource) Allow(method string) *Resource {
	res.methods[strings.ToUpper(method)] = struct{}{}
	return res
}

func (res *Resource) Disallow(method string) *Resource {
	delete(res.methods, strings.ToUpper(method))
	return res
}

func newResource(name string, structType interface{}, api *API) *Resource {
	t := reflect.TypeOf(structType)

	// If ptr, dereference schema type
	if t.Kind().String() == "ptr" {
		t = t.Elem()
	}

	res := Resource{Name: name, Type: t, api: api}

	// initialize the allowed methods
	res.methods = make(map[string]struct{})
	res.methods["GET"] = struct{}{}
	res.methods["POST"] = struct{}{}
	res.methods["DELETE"] = struct{}{}

	res.ItemHandlers = make(map[string]func(http.ResponseWriter, *http.Request, interface{}))
	res.ItemHandlers["GET"] = defaultItemGET
	res.ItemHandlers["DELETE"] = defaultItemDELETE

	res.CollectionHandlers = make(map[string]func(http.ResponseWriter, *http.Request, []interface{}))
	res.CollectionHandlers["GET"] = defaultCollectionGET
	res.CollectionHandlers["POST"] = defaultCollectionPOST

	return &res
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

	// Check if method is allowed and
	// handler has been implemented for the request method
	_, allowed := res.methods[r.Method]
	_, exists := res.ItemHandlers[r.Method]
	if exists && allowed {
		// query database
		var obj interface{}
		if r.Method == "DELETE" {
			obj = api.DB.DeleteById(res.Type, res.Name, id)
		} else {
			obj = api.DB.FindById(res.Type, res.Name, id)
		}

		// call handler
		res.ItemHandlers[r.Method](w, r, obj)
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

	// Check if method is allowed and
	// handler has been implemented for the request method
	_, allowed := res.methods[r.Method]
	_, exists := res.CollectionHandlers[r.Method]
	if exists && allowed {
		var objs []interface{}
		if r.Method == "POST" {
			obj := reflect.New(res.Type).Interface()

			// limit body size to avoid injection
			body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
			if err != nil {
				panic(err)
			}
			if err := r.Body.Close(); err != nil {
				panic(err)
			}
			if err := json.Unmarshal(body, &obj); err != nil {
				panic(err)
			}

			// insert into database
			objs = api.DB.Insert(obj, res.Name)
		} else {
			// read from database
			objs = api.DB.FindAll(res.Type, res.Name)
		}
		res.CollectionHandlers[r.Method](w, r, objs)
	} else {
		fmt.Println(r.RequestURI + " does not have a " + r.Method + " method defined")
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
	}

	return
}
