package api

import (
	"net/http"
	"reflect"
	//	"../mux"
	"encoding/json"
)

type Resource struct {
	Name     string
	Type     reflect.Type
	Handlers map[string]func(interface{}, http.ResponseWriter, []string)
}

func NewResource(name string, structType interface{}) *Resource {
	t := reflect.TypeOf(structType)

	// If ptr, dereference schema type
	if t.Kind().String() == "ptr" {
		t = t.Elem()
	}

	r := Resource{Name: name, Type: t}
	r.Handlers = make(map[string]func(interface{}, http.ResponseWriter, []string))
	r.Handlers["GET"] = defaultGET
	r.Handlers["PUT"] = defaultPUT

	return &r
}

func defaultGET(obj interface{}, w http.ResponseWriter, accepts []string) {
	if len(accepts) > 1 {
		panic("ERROR with GET: Override the GET handler to support accepts other than application/json")
	} else if accepts[0] != "application/json" {
		panic("ERROR with GET: Override the GET handler to support accepts other than application/json")
	}

	// if object was found in the database
	if obj != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(obj)
		if err != nil {
			panic(err)
		}
	} else {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		jsonErr := struct {
			Code int
			Text string
		}{Code: http.StatusNotFound, Text: "Not Found"}
		err := json.NewEncoder(w).Encode(jsonErr)
		if err != nil {
			panic(err)
		}
	}

}

func defaultPUT(interface{}, http.ResponseWriter, []string) {
}

type DefaultGETHandler struct {
	res *Resource
}

func (h DefaultGETHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// find id from request
	//	id := mux.GetParams(r)["id"]

	// empty struct -- TODO: change to read from database
	obj := reflect.New(h.res.Type).Elem().Interface()

	h.res.Handlers["GET"](obj, w, []string{"application/json"})

	//fmt.Fprintf(w, "<h1>You've reached resource " + h.res.Name + " with id " + id + "!</h1>")
	return
}
