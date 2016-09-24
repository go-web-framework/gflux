//Need to decide on proper package arrangement later
//For now edit import to proper mux package
package main

import (
	"fmt"
	"net/http"
	"./mux"
	"./bs_radix"
)

func main(){
	testMux := mux.New()
	testHandler1 := testHandler{}
	testMux.Handle("/test", nil, testHandler1)
	testRadix := bs_radix.New()
	testRadix.Insert("cat", 1)
	testRadix.Insert("catch", 2)
	testRadix.Insert("catcher", 2)
	fmt.Println("testRadix is ", testRadix)
	testRadix.Insert("bob", 50)
	testRadix.Insert("catty", 6)
	testRadix.Insert("catter", 10)
	fmt.Println("size of tree is now ", testRadix.Len())
	fmt.Println("Now listening at :8080 ... ")
	http.ListenAndServe(":8080", testMux)


}

type testHandler struct{
}

func (t testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "<h1>Test</h1>")
	return
}

