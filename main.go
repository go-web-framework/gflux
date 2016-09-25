package main

import (
	"fmt"
	"net/http"
	"./mux"
)

func main(){
	testMux := mux.New()
	testHandler1 := testHandler{}
	testHandler2 := homeHandler{}
	testHandler3 := handler404{}
	testMux.Handle("/test", nil, testHandler1)
	testMux.Handle("/home", nil, testHandler2)
	testMux.SetNotFound(testHandler3)
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

type handler404 struct{
}

func (t handler404) ServeHTTP(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "<h1>You've reached a custom 404!</h1>")
	return
}
