package main

import (
	"./api"
)

type Post struct {
	Title string
	Author string
}

func main() {
	a := api.New("test.db")
	defer a.Close()
	
	a.NewResource("posts", &Post{})
	
	a.Serve()
}
