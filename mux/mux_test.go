//Need to decide on proper package arrangement later
//For now edit import to proper mux package
package mux

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

//assumes request is accessed for GetParams
func checkParams(r *http.Request ,keys []string ,expected []string ) (bool, string){
	if (len(keys) != len(expected)){
		return false, "Length keys and answers do not match"
	}
	nullKey := "Key null: "
	wrongVal := "Wrong value: "
	params := GetParams(r);
	for i := 0; i < len(keys); i++ {
		a := params[keys[i]]
		if (a == ""){
			return false, nullKey + keys[i]
		}
		if (a != expected[i]){
			return false, wrongVal + keys[i] + " " + expected[i]
		}
	}
	return true, ""
}

// Wildcards as every fragment of the path
func TestWildcard2(t *testing.T) {
	mux := New()
	call := false
	paramCheck := false
	errStr := ""
	mux.Handle("/{a}/{b}/{c}", nil, http.HandlerFunc(func(w http.ResponseWriter,r2 *http.Request) {
		call = true
		paramCheck, errStr = checkParams(r2, []string{"a", "b", "c"}, []string{"e", "f", "g"})
	}))

	r, _ := http.NewRequest("POST", "/e/f/g", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	if !call {
		t.Error("handler should be called")
	}
	if !paramCheck{
		t.Error(errStr)
	}
}

// Wildcards in various places
func TestWildcard3(t *testing.T) {
	mux := New()
	call := false
	paramCheck := false
	errStr := ""
	mux.Handle("/a/", nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		call = false
	}))
	mux.Handle("/B/b", nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		call = false
	}))

	mux.Handle("/ap/{index}",
		nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
			call = false
		}))

	mux.Handle("/ap/socialinjustice/jane/{letter}",
		nil, http.HandlerFunc(func(w http.ResponseWriter, r2 *http.Request) {
			call = true
			paramCheck, errStr = checkParams(r2, []string{"letter"}, []string{"delta"})
		}))

	mux.Handle("/ap/socialwarriors/avenged/",
		nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
			call = false
		}))

	r, _ := http.NewRequest("GET", "/ap/socialinjustice/jane/delta", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	if !call {
		t.Error("handler should be called")
	}
	if !paramCheck{
		t.Error(errStr)
	}
}

// Wildcards in various places
func TestWildcard4(t *testing.T) {
	mux := New()
	call := false
	paramCheck := false
	errStr := ""
	mux.Handle("/a/", nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		call = false
	}))
	mux.Handle("/B/b", nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		call = false
	}))

	mux.Handle("/ap/{index}",
		nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
			call = false
		}))

	mux.Handle("/ap/socialinjustice/jane/{letter}",
		nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
			call = false
		}))

	mux.Handle("/ap/socialinjustice/{type}/delta",
		nil, http.HandlerFunc(func(w http.ResponseWriter, r2 *http.Request) {
			call = true
			paramCheck, errStr = checkParams(r2, []string{"type"}, []string{"november"})
		}))

	mux.Handle("/ap/socialwarriors/avenged/",
		nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
			call = false
		}))

	r, _ := http.NewRequest("GET", "/ap/socialinjustice/november/delta", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	if !call {
		t.Error("handler should be called")
	}
	if !paramCheck{
		t.Error(errStr)
	}

}

// Test found wildcard with different tree structure
func TestWildcard5(t *testing.T) {
	mux := New()
	call := false
	paramCheck := false
	errStr := ""
	mux.Handle("/a/", nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		call = false
	}))

	mux.Handle("/ap/socialinjustice/{index}/delta",
		nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
			call = false
		}))

	mux.Handle("/ap/socialinjustice/{index}/delta/{month}",
		nil, http.HandlerFunc(func(w http.ResponseWriter, r2 *http.Request) {
			call = true
			paramCheck, errStr = checkParams(r2, []string{"index", "month"}, []string{"november", "nope"})
		}))

	mux.Handle("/ap/socialwarriors/avenged/",
		nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
			call = false
		}))

	r, _ := http.NewRequest("GET", "/ap/socialinjustice/november/delta/nope", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	if !call {
		t.Error("handler should be called")
	}
	if !paramCheck{
		t.Error(errStr)
	}
}

// Test found wildcard at end of longer Url path
func TestRoutes(t *testing.T) {
	mux := New()
	call := false
	mux.Handle("/a/", nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		call = false
	}))

	mux.Handle("/ag/a",
		nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
			call = true
		}))

	r, _ := http.NewRequest("GET", "/ag/a", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	if !call {
		t.Error("handler should be called")
	}
}

// Test found wildcard at end of longer Url path
func TestHome(t *testing.T) {
	mux := New()
	call := false
	mux.Handle("/", nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		call = true
	}))

	mux.Handle("/ag/a",
		nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
			call = false
		}))

	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	if !call {
		t.Error("handler should be called")
	}
}

// Test the custom not handler handler sets 404 error code
func TestNotFoundCustomHandlerSends404(t *testing.T) {
	mux := New()
	mux.NotFound = http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(404)
		rw.Write([]byte("You've reached a custom 404!"))
	})

	r, _ := http.NewRequest("GET", "/b/123", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	if w.Code != 404 {
		t.Errorf("expecting error code 404, got %v", w.Code)
	}
}

// Test the custom not handler handler sets 404 error code
func TestNotFoundDefaultHandlerSends404(t *testing.T) {
	mux := New()

	r, _ := http.NewRequest("GET", "/b/123", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	if w.Code != 404 {
		t.Errorf("expecting error code 404, got %v", w.Code)
	}
}

// Test forward slash compatibility
func TestForwardSlashBehavior1(t *testing.T) {
	mux := New()
	call := false
	mux.Handle("/a", nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		call = true
	}))

	r, _ := http.NewRequest("GET", "/a/", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	if !call {
		t.Error("handler should be called")
	}
}

// Test forward slash compatibility
func TestForwardSlashBehavior2(t *testing.T) {
	mux := New()
	call := false
	mux.Handle("/a/", nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		call = true
	}))

	r, _ := http.NewRequest("GET", "/a", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	if !call {
		t.Error("handler should be called")
	}
}

// Test if path not found in entirety
func TestNotFoundPathBeginning(t *testing.T) {
	mux := New()
	call := false
	mux.Handle("/a/", nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		call = false
	}))
	mux.Handle("/a/b", nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		call = false
	}))

	mux.Handle("/a/c",
		nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
			call = true
		}))

	r, _ := http.NewRequest("GET", "/d", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	if call {
		t.Error("handler should not be called")
	}
}

// Ensure immediate paths still found
func TestFoundPathBeginning(t *testing.T) {
	mux := New()
	call := false
	mux.Handle("/a", nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		call = true
	}))
	mux.Handle("/a/b", nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		call = false
	}))

	mux.Handle("/a/c",
		nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
			call = false
		}))

	r, _ := http.NewRequest("GET", "/a", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	if !call {
		t.Error("handler should be called")
	}
}

// ensure URI fragments do not, by themselves, 
// do not cause an attempt to return a handler if requested
func TestTreeStructure(t *testing.T) {
	mux := New()
	call := false
	mux.Handle("/a/", nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		call = true
	}))
	mux.Handle("/B/b", nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		call = true
	}))

	r, _ := http.NewRequest("GET", "/B", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	if call {
		t.Error("handler should not be called")
	}
}

// Using the walk function of the radix tree, ensuring
// proper return
func TestTreeStructure2(t *testing.T) {
	mux := New()
	call := false
	mux.Handle("/a/", nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		call = false
	}))

	mux.Handle("/ap",
		nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
			call = true
		}))

	found := mux.trie.Walk("/a")

	if !found {
		t.Error("Walk should return true")
	}
}

// Test found wildcard with different tree structure
func TestRouting8(t *testing.T) {
	mux := New()
	call := false
	mux.Handle("/a/{b}/c/d", nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		call = false
	}))

	defer func() {
		if p := recover(); p == nil {
			t.Error("route registration should panic")
		}
	}()

	mux.Handle("/a/{b}/c/d", nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		call = false
	}))
}

// Test found wildcard with different tree structure
func TestRouting9(t *testing.T) {
	mux := New()
	call := false
	route := mux.Handle("/a/{b}/c/d", nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		call = false
	}))

	if route == nil {
		t.Error("route should not be nil")
	}

	defer func() {
		if p := recover(); p == nil {
			t.Error("route registration should panic")
		}
	}()

	mux.Handle("/a/{a}/c/d", nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		call = false
	}))
}

// Test adding a value to an existing, non-leaf fragment
func TestRouting10(t *testing.T) {
	mux := New()
	call := false
	mux.Handle("/a/b/", nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		call = false
	}))
	mux.Handle("/a/", nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		call = true
	}))
	mux.Handle("/a/c", nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		call = true
	}))

	r, _ := http.NewRequest("GET", "/a", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	if !call {
		t.Error("handler should be called")
	}
}


// Attempt to directly request a wildcard - should not return the wildcard handler
func TestDirectAccess(t *testing.T) {
	mux := New()
	call := false
	mux.Handle("/a/{b}/c", nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		call = true
	}))

	r, _ := http.NewRequest("GET", "/a/{b}/c", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	if call {
		t.Error("handler should not be called")
	}
}

// Create a GET handler for a path. GET request should return the handler
func TestMethodHandling1(t *testing.T) {
	mux := New()
	call := false
	mux.GET("/a", nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		call = true
	}))

	r, _ := http.NewRequest("GET", "/a/", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	if !call {
		t.Error("handler should be called")
	}
}

// Create a GET handler for a path. POST request should not return the handler
func TestMethodHandling2(t *testing.T) {
	mux := New()
	call := false
	mux.GET("/a", nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		call = true
	}))

	r, _ := http.NewRequest("POST", "/a/", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	if call {
		t.Error("handler should not be called")
	}
}

// Test for a panic when attempting to register the same path
// multiple times with the same method
func TestMethodHandling3(t *testing.T) {
	mux := New()
	call := false
	mux.GET("/a", nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		call = true
	}))

	defer func() {
		if p := recover(); p == nil {
			t.Error("route registration should panic")
		}
	}()

	mux.GET("/a", nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		call = true
	}))
}

// Create a POST handler only for a path. POST request should return the handler
func TestMethodHandling4(t *testing.T) {
	mux := New()
	call := true

	mux.GET("/a", nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		call = false
	}))
	mux.POST("/a", nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		call = false
	}))
	mux.DELETE("/a", nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		call = false
	}))
	mux.HEAD("/a", nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		call = false
	}))
	mux.PUT("/a", nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		call = false
	}))
	mux.OPTIONS("/a", nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		call = false
	}))
	mux.PATCH("/a", nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		call = true
	}))
	mux.Handle("/a", nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		call = true
	}))

	r, _ := http.NewRequest("POST", "/a", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	if !call {
		t.Error("handler should be called")
	}
}

// Test method handling with wildcards 
// create many methods for a wildcard path
// ensuring the proper path is still returned
func TestWildcardMethodHandling(t *testing.T) {
		mux := New()
	call := false
	paramCheck := false
	errStr := ""
	mux.GET("/{a}/{b}/{c}", nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		call = false
	}))
	mux.POST("/{a}/{b}/{c}", nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		call = false
	}))
	mux.DELETE("/{a}/{b}/{c}", nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		call = false
	}))
	mux.HEAD("/{a}/{b}/{c}", nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		call = false
	}))
	mux.PUT("/{a}/{b}/{c}", nil, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		call = true
		paramCheck, errStr = checkParams(req, []string{"a", "b", "c"}, []string{"e", "f", "g"})
	}))
	mux.OPTIONS("/{a}/{b}/{c}", nil, http.HandlerFunc(func(w http.ResponseWriter, r2 *http.Request) {
		call = false
	}))
	mux.PATCH("/{a}/{b}/{c}", nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		call = false
	}))


	r, _ := http.NewRequest("PUT", "/e/f/g/", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	if !call {
		t.Error("handler should be called")
	}
	if !paramCheck{
		t.Error(errStr)
	}
}
