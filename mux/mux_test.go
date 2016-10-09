//Need to decide on proper package arrangement later
//For now edit import to proper mux package
package mux

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test found wildcard at end of Url path
func TestWildcard2(t *testing.T) {
	mux := New()
	call := false

	
	mux.Handle("/{a}/{b}/{c}", nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		call = true
	}))

	r, _ := http.NewRequest("GET", "/e/f/g", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	if !call {
		t.Error("handler should be called")
	}
}

// Test found wildcard with different tree structure
func TestWildcard3(t *testing.T) {
	mux := New()
	call := false
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
		call = true
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
}


// Test found wildcard with different tree structure
func TestWildcard4(t *testing.T) {
	mux := New()
	call := false
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
		nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		call = true
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

}

// Test found wildcard with different tree structure
func TestWildcard5(t *testing.T) {
	mux := New()
	call := false
	mux.Handle("/a/", nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		call = false
	}))

	mux.Handle("/ap/socialinjustice/{index}/delta", 
		nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		call = false
	}))

	 mux.Handle("/ap/socialinjustice/{index}/delta/{month}", 
		nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		call = true
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
	mux.SetNotFound(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(404)
		rw.Write([]byte("You've reached a custom 404!"))
	}))

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



// Test not found
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

// Test not found
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


// Test found wildcard at end of longer Url path
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

// Test found wildcard at end of longer Url path
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

	found := mux.radix.Walk("/a")

	if !found{
		t.Error("handler should be called")
	}
}

// Test found wildcard with different tree structure
func TestRouting8(t *testing.T) {
	mux := New()
	call := false
	route := mux.Handle("/a/{b}/c/d", nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		call = false
	}))
	
	if route == nil {
		t.Error("route should not be nil")
	}

	route2 := mux.Handle("/a/{b}/c/d", nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		call = false
	}))

	if route2 != nil {
		t.Error("route should be nil")
	}
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

	route2 := mux.Handle("/a/{a}/c/d", nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		call = false
	}))

	if route2 != nil {
		t.Error("route should be nil")
	}
}

// Test found wildcard with different tree structure
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

