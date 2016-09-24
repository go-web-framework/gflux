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
	/*
	testRadix.Insert("/home", 2)
	fmt.Println("size of tree is now ", testRadix.Len())
	//fmt.Println("Now listening at :8080 ... ")
	//http.ListenAndServe(":8080", testMux)
	error := 0

	//BEGIN TEST : ensure both "/SDP" and "/SDP/" return the same value
	fmt.Println("")
	error = 0
	testRadix.Insert("/SDP", 1234)
	size1 := testRadix.Len()
	testRadix.Insert("/SDP/", 1234)
	val, found := testRadix.Get("/SDP")
	val2, found2 := testRadix.Get("/SDP/")
	size2 := testRadix.Len()
	
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
		fmt.Println("/SDP and /SDP/ test successful")
	}
	//END TEST

	//BEGIN TEST : wildcard placeholder test
	error = 0
	fmt.Println("")
	fmt.Println("Starting size of tree : ", testRadix.Len())
	fmt.Println("Inserting (/EE445L, 400) ")
	fmt.Println("Inserting (/EE445L/Aug, 500) ")
	fmt.Println("Inserting (/EE445L/{:Id}, 600) ")
	testRadix.Insert("/EE445L", 400)
	testRadix.Insert("/EE445L/Aug", 500)
	testRadix.Insert("/EE445L/{:Id}", 600)
	fmt.Println("Size of tree is now ", testRadix.Len())
	fmt.Println("Attempting to search for /EE445L/Sept should return value of 600")
	val, found = testRadix.Get("/EE445L/Sept")
	if found {
		fmt.Println("/EE445L/Sept was found. Val : ", val)
	} else {
			fmt.Println("/EE445L/Sept was not found. ")
			error ++
	}

	if error == 0 {
		fmt.Println("/{:Id} test successful")
	}
	//END TEST
*/
	fmt.Println("Inserting (/EE445L, 400) ")
	testRadix.Insert("/EE445L", 400)
	testRadix.Insert("/EE445L/Aug", 500)
	testRadix.Insert("/EE445L/Sept", 600)
	testRadix.Insert("/EE445L/Sands", 700)
	testRadix.Insert("/EE445L/{:Id}", 800)
	_, found := testRadix.Get("/EE445L/Sands")
	if found {
		fmt.Println("********************/EE445L/Sands was found")
	}
	_, found = testRadix.Get("/EE445L/Oct")
	if found {
		fmt.Println("**********************/EE445L/Oct was found")
	}
}

type testHandler struct{
}

func (t testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "<h1>Test</h1>")
	return
}

