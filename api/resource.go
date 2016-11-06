package api

import (
	"net/http"
	"reflect"
)

type Resource struct {
	Name string
	Type reflect.Type
	Handlers map[string]func(interface{}, http.ResponseWriter, []string)
}

func NewResource(name string, t reflect.Type) *Resource{
	r := Resource{Name: name, Type: t}
	r.Handlers = make(map[string]func(interface{}, http.ResponseWriter, []string))
	r.Handlers["GET"] = defaultGET
	r.Handlers["PUT"] = defaultPUT
	
	return &r
}

func defaultGET(interface{}, http.ResponseWriter, []string) {
}

func defaultPUT(interface{}, http.ResponseWriter, []string) {
}
