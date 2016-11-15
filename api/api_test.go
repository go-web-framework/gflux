package api

import "testing"

type Post struct {
	Title  string
	Author string
}

func ExampleAPI_NewResource() {
	a := New("sqlite3", "test.db")
	defer a.Close()

	a.NewResource("posts", &Post{})

	a.Serve()
}

func TestAPI_NewResource_NonPointer(t *testing.T) {
	a := New("sqlite3", "test.db")
	defer a.Close()

	a.NewResource("posts", Post{}) // Post{} instead of &Post{}
}
