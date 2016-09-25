//Need to decide on proper package arrangement later
//For now edit import to proper mux package
package bs_radix

import (
	"net/http"
	"testing"
	"github.com/stretchr/testify/assert"
)


//ensure both "/SDP" and "/SDP/" return the same value
func TestForwardSlashBehavior( t *testing.T){
	testRadix := New()
	testHandler := testHandler_t{}
	testRadix.NewRoute("/SDP", testHandler)
	size1 := testRadix.Len()
	assert.Equal(t, size1, 1, "Size is not 1 after first insert")
	testRadix.NewRoute("/SDP/", testHandler)
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
	testRadix := New()
	testHandler := testHandler_t{}
	testRadix.NewRoute("/EE445L", testHandler)
	testRadix.NewRoute("/EE445L/Aug", testHandler)
	testRadix.NewRoute("/EE445L/Sept", testHandler)
	testRadix.NewRoute("/EE445L/Sands", testHandler)
	testRadix.NewRoute("/EE445L/{:Id}", testHandler)
	assert.Equal(t, testRadix.Len(), 5)
	_, found := testRadix.Get("/EE445L/Sands")
	assert.True(t, found)
	_, found = testRadix.Get("/EE445L/Oct")
	assert.True(t, found)
}


type testHandler_t struct{
}
func (t testHandler_t) ServeHTTP(w http.ResponseWriter, r *http.Request){
	
	return
}



