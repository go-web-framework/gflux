//Need to decide on proper package arrangement later
//For now edit import to proper mux package
package main

import (
	"fmt"
	"net/http"
	"./mux"
	"testing"
)

func TestMux(t *testing.T){
	
	testMux := mux.New()
	testHandler1 := testHandler{}
	testHandler2 := homeHandler{}
	testMux.Handle("/test", nil, testHandler1)
	testMux.Handle("/home", nil, testHandler2)
	fmt.Println("Listening on :8080")
	http.ListenAndServe(":8080", testMux)
	
}


type testHandler struct{
}

func (t testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "<h1>Test</h1>")
	return
}

type homeHandler struct{
}

func (t homeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "<h1>Home</h1>")
	return
}

