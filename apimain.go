package main

import (
	"./api"
)

type Post struct {
	api.Model
	Title string
	Author string
}

func main() {
	a, err := api.New("test.db")
	if err != nil {
		panic("failed to connect to database")
	}
	defer a.Close()
	
	a.NewResource("posts", &Post{})
	
	a.Serve()
}
