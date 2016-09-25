//Need to decide on proper package arrangement later
//For now edit import to proper mux package
package main

import (
	"fmt"
	"net/http"
	"./mux"
	"./bs_radix"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestMux(t *testing.T){
	
	testMux := mux.New()
	testHandler1 := testHandler{}
	testMux.Handle("/test", nil, testHandler1)
	/*
	
	//http.ListenAndServe(":8080", testMux)
	*/
}

//ensure both "/SDP" and "/SDP/" return the same value
func TestForwardSlashBehavior( t *testing.T){
	testRadix := bs_radix.New()
	testRadix.Insert("/SDP", 1234)
	size1 := testRadix.Len()
	assert.Equal(t, size1, 1, "Size is not 1 after first insert")
	testRadix.Insert("/SDP/", 1234)
	val1, found := testRadix.Get("/SDP")
	val2, found2 := testRadix.Get("/SDP/")
	size2 := testRadix.Len()
	
	assert.Equal(t, size1, size2, "Size should not change")
	assert.True(t, found)
	assert.True(t, found2)
	assert.Equal(t, val1, val2)
}

//Test in progress: working to ensure that wildcards act as expected
func TestWildcard(t *testing.T) {
	testRadix := bs_radix.New()
	testRadix.Insert("/EE445L", 400)
	testRadix.Insert("/EE445L/Aug", 500)
	testRadix.Insert("/EE445L/Sept", 600)
	testRadix.Insert("/EE445L/Sands", 700)
	testRadix.Insert("/EE445L/{:Id}", 800)
	assert.Equal(t, testRadix.Len(), 5)
	_, found := testRadix.Get("/EE445L/Sands")
	assert.True(t, found)
	_, found = testRadix.Get("/EE445L/Oct")
	assert.True(t, found)
}


type testHandler struct{
}

func (t testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "<h1>Test</h1>")
	return
}

