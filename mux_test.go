//Need to decide on proper package arrangement later
//For now edit import to proper mux package
package main

import (
	"net/http"
	"./mux"
	"net/http/httptest"
	"testing"
)

// Test if the route is valid
func TestRouting1(t *testing.T) {
	mux := mux.New()
	call := false
	mux.Handle("/a/", nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		call = true
	}))

	r, _ := http.NewRequest("GET", "/b/123", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	if call {
		t.Error("handler should not be called")
	}
}

// Test if the route is valid
func TestRouting2(t *testing.T) {
	mux := mux.New()
	call := false
	mux.Handle("/a/", nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		call = true
	}))

	r, _ := http.NewRequest("GET", "/a/", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	if !call {
		t.Error("handler should be called")
	}
}

// Test the custom not handler handler sets 404 error code
func TestNotFoundCustomHandlerSends404(t *testing.T) {
	mux := mux.New()
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
	mux := mux.New()

	r, _ := http.NewRequest("GET", "/b/123", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	if w.Code != 404 {
		t.Errorf("expecting error code 404, got %v", w.Code)
	}
}

// Test forward slash compatibility
func TestForwardSlashBehavior1(t *testing.T) {
	mux := mux.New()
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
	mux := mux.New()
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

// Test found wildcard at end of Url path
func TestWildcard1(t *testing.T) {
	mux := mux.New()
	call := false
	mux.Handle("/a/", nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		call = false
	}))
	mux.Handle("/a/b", nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		call = false
	}))
	mux.Handle("/a/{:Id}", nil, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		call = true
	}))

	r, _ := http.NewRequest("GET", "/a/c", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	if !call {
		t.Error("handler should be called")
	}
}

