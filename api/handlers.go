package api

import (
	"encoding/json"
	"net/http"
	"strings"
)

func defaultItemGET(w http.ResponseWriter, r *http.Request, obj interface{}) {
	accept := r.Header.Get("Accept")
	if !strings.Contains(accept, "application/json") && !strings.Contains(accept, "*/*") {
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

func defaultCollectionGET(w http.ResponseWriter, r *http.Request, objs []interface{}) {
	accept := r.Header.Get("Accept")
	if !strings.Contains(accept, "application/json") && !strings.Contains(accept, "*/*") {
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

func defaultItemDELETE(w http.ResponseWriter, r *http.Request, obj interface{}) {
	accept := r.Header.Get("Accept")
	if !strings.Contains(accept, "application/json") && !strings.Contains(accept, "*/*") {
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

func defaultCollectionPOST(w http.ResponseWriter, r *http.Request, objs []interface{}) {
	accept := r.Header.Get("Accept")
	if !strings.Contains(accept, "application/json") && !strings.Contains(accept, "*/*") {
		panic("ERROR with POST: Override the POST CollectionHandler to support accepts other than application/json")
	}

	// if object was put into database
	if objs != nil && len(objs) > 0 && objs[0] != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		err := json.NewEncoder(w).Encode(objs[0])
		if err != nil {
			panic(err)
		}
	} else {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusUnprocessableEntity)
		jsonErr := struct {
			Code int
			Text string
		}{Code: http.StatusUnprocessableEntity, Text: "Unprocessable Entity"}
		err := json.NewEncoder(w).Encode(jsonErr)
		if err != nil {
			panic(err)
		}
	}
}
