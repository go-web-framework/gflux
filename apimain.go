package main

import (
    "./api"
)

type Post struct {
    Title string
    Author string
}

func main() {
    a := api.New("sqlite3", "test.db")
    //	a := api.New("mysql", "user:password@/test")
    defer a.Close()

    a.NewResource("posts", &Post{})

    a.Serve()
}
