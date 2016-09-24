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
	testRadix.Insert("/home", 2)
	testRadix.Insert("/ee445L/pg", 2)
	fmt.Println("testRadix is ", testRadix)
	fmt.Println("size of tree is now ", testRadix.Len())
	fmt.Println("size of tree is now ", testRadix.Len())
	//fmt.Println("Now listening at :8080 ... ")
	//http.ListenAndServe(":8080", testMux)

	//BEGIN TEST : ensure both "/SDP" and "/SDP/" return the same value
	testRadix.Insert("/SDP", 1234)
	size1 := testRadix.Len()
	testRadix.Insert("/SDP/", 1234)
	val, found := testRadix.Get("/SDP")
	val2, found2 := testRadix.Get("/SDP/")
	size2 := testRadix.Len()
	error := 0
	if size1 != size2 {
		fmt.Println("Size of tree remained same")
		error++;
	}
	if !found || !found2 {
		fmt.Println("Did not find both /SDP & /SDP/")
		error++
	}
	if val != val2 {
		fmt.Println("Both were found, but the values were different")
		error++
	}
	if error == 0 {
		fmt.Println("Attempting to insert /SDP and /SDP/ did not cause multiple entries")
		fmt.Println("Both /SDP and /SDP/ are both returned by Get call and return the same value")
		fmt.Println("/SDP and /SDP/ test successful")

	}
	//END TEST
}

type testHandler struct{
}

func (t testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "<h1>Test</h1>")
	return
}

