// package api allows for easy creation of REST APIs
package api

type Post struct {
    Title string
    Author string
}

func ExampleAPI_NewResource() {
    a := New("sqlite3", "test.db")
    defer a.Close()

    a.NewResource("posts", &Post{})

    a.Serve()
}
