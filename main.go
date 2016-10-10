package main

import (
	"fmt"
	"net/http"

	"github.com/go-web-framework/gflux/mux"
)

func main() {
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

type testHandler struct {
}

func (t testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//myVars := r.Context().Value(int32(0)).(map[string]string)
	// a := mux.Vars(r)
	//a := myVars["Var1"]
	//b := myVars["2"]
	// fmt.Fprintf(w, "<h1>Test %s</h1>", a)
	return
}

type homeHandler struct {
}

func (t homeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Home</h1>")
	return
}

type handler404 struct {
}

func (t handler404) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>You've reached a custom 404!</h1>")
	return
}
