// package api allows for easy creation of REST APIs
package api

import (
	"encoding/json"
	"net/http"
)

func defaultItemGET(obj interface{}, w http.ResponseWriter, accepts []string) {
	if len(accepts) > 1 {
		panic("ERROR with GET: Override the GET ItemHandler to support accepts other than application/json")
	} else if accepts[0] != "application/json" {
		panic("ERROR with GET: Override the GET ItemHandler to support accepts other than application/json")
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

func defaultCollectionGET(objs []interface{}, w http.ResponseWriter, accepts []string) {
	if len(accepts) > 1 {
		panic("ERROR with GET: Override the GET CollectionHandler to support accepts other than application/json")
	} else if accepts[0] != "application/json" {
		panic("ERROR with GET: Override the GET CollectionHandler to support accepts other than application/json")
	}

	// if objects were found in the database
	if objs != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(objs)
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

func defaultItemDELETE(obj interface{}, w http.ResponseWriter, accepts []string) {
	if len(accepts) > 1 {
		panic("ERROR with DELETE: Override the DELETE ItemHandler to support accepts other than application/json")
	} else if accepts[0] != "application/json" {
		panic("ERROR with DELETE: Override the DELETE ItemHandler to support accepts other than application/json")
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

//func defaultCollectionDELETE(objs []interface{}, w http.ResponseWriter, accepts []string) {
//	if len(accepts) > 1 {
//		panic("ERROR with DELETE: Override the DELETE CollectionHandler to support accepts other than application/json")
//	} else if accepts[0] != "application/json" {
//		panic("ERROR with DELETE: Override the DELETE CollectionHandler to support accepts other than application/json")
//	}

//	// always return 404
//	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
//	w.WriteHeader(http.StatusNotFound)
//	jsonErr := struct {
//		Code int
//		Text string
//	}{Code: http.StatusNotFound, Text: "Not Found"}
//	err := json.NewEncoder(w).Encode(jsonErr)
//	if err != nil {
//		panic(err)
//	}
//}

func defaultItemPUT(interface{}, http.ResponseWriter, []string) {
}
